var builder = WebApplication.CreateBuilder(args);

builder.Services.AddEndpointsApiExplorer(); 
builder.Services.AddSwaggerGen(c =>
{
    c.AddSecurityDefinition("ApiKey", new Microsoft.OpenApi.Models.OpenApiSecurityScheme
    {
        Description = "API Key needed to access the endpoints. X-API-Key: {API Key}",
        In = Microsoft.OpenApi.Models.ParameterLocation.Header,
        Name = "X-API-Key",
        Type = Microsoft.OpenApi.Models.SecuritySchemeType.ApiKey,
        Scheme = "ApiKey"
    });

    c.AddSecurityRequirement(new Microsoft.OpenApi.Models.OpenApiSecurityRequirement
    {
        {
            new Microsoft.OpenApi.Models.OpenApiSecurityScheme
            {
                Reference = new Microsoft.OpenApi.Models.OpenApiReference
                {
                    Type = Microsoft.OpenApi.Models.ReferenceType.SecurityScheme,
                    Id = "ApiKey"
                }
            },
            new string[] {}
        }
    });
});

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();

var apiKey = builder.Configuration["ApiKey"];

app.Use(async (context, next) =>
{
    if (!context.Request.Headers.TryGetValue("X-API-Key", out var extractedApiKey) ||
        extractedApiKey != apiKey)
    {
        context.Response.StatusCode = 401;
        await context.Response.WriteAsync("Unauthorized");
        return;
    }
    await next();
});

try
{
    var connString = builder.Configuration.GetConnectionString("DefaultConnection");
    if (!string.IsNullOrEmpty(connString))
    {
        using var conn = new Npgsql.NpgsqlConnection(connString);
        conn.Open();
        app.Logger.LogInformation("Successfully connected to PostgreSQL");
        conn.Close();
    }
    else
    {
        app.Logger.LogWarning("No DefaultConnection string found in configuration.");
    }
}
catch (Exception ex)
{
    app.Logger.LogError(ex, "Failed to connect to PostgreSQL on startup");
}

app.MapGet("/okCode", () =>
{
    return Results.Ok("Everything is awesome!");
}).WithName("getOkCode")
.WithOpenApi();

app.MapGet("/users", async (HttpContext context) =>
{
    var connString = app.Configuration.GetConnectionString("DefaultConnection");
    if (string.IsNullOrWhiteSpace(connString))
    {
        app.Logger.LogWarning("Attempt to call /users but DefaultConnection is not configured.");
        return Results.Problem(detail: "Database connection is not configured.", statusCode: 500);
    }

    try
    {
        await using var conn = new Npgsql.NpgsqlConnection(connString);
        await conn.OpenAsync();
        await using var cmd = new Npgsql.NpgsqlCommand("SELECT * FROM users", conn);
        await using var reader = await cmd.ExecuteReaderAsync();

        var results = new List<Dictionary<string, object?>>();
        while (await reader.ReadAsync())
        {
            var row = new Dictionary<string, object?>();
            for (int i = 0; i < reader.FieldCount; i++)
            {
                var name = reader.GetName(i);
                var value = await reader.IsDBNullAsync(i) ? null : reader.GetValue(i);
                row[name] = value;
            }
            results.Add(row);
        }

        return Results.Ok(results);
    }
    catch (Exception ex)
    {
        app.Logger.LogError(ex, "Error while fetching users from database");
        return Results.Problem(detail: "An error occurred while fetching users.", statusCode: 500);
    }
}).WithName("getUsers").WithOpenApi();

app.MapPost("/users", async (UserCreateDto dto) =>
{
    if (dto == null)
    {
        return Results.BadRequest(new { error = "Request body is required." });
    }

    if (string.IsNullOrWhiteSpace(dto.Name) || string.IsNullOrWhiteSpace(dto.Username) || string.IsNullOrWhiteSpace(dto.Password))
    {
        return Results.BadRequest(new { error = "Name, Username and Password are required to create a user." });
    }

    var connString = app.Configuration.GetConnectionString("DefaultConnection");
    if (string.IsNullOrWhiteSpace(connString))
    {
        app.Logger.LogWarning("Attempt to call POST /users but DefaultConnection is not configured.");
        return Results.Problem(detail: "Database connection is not configured.", statusCode: 500);
    }

    try
    {
        await using var conn = new Npgsql.NpgsqlConnection(connString);
        await conn.OpenAsync();

    const string sql = @"INSERT INTO users (name, username, password) VALUES (@name, @username, @password) RETURNING id";

        await using var cmd = new Npgsql.NpgsqlCommand(sql, conn);
        cmd.Parameters.AddWithValue("name", NpgsqlTypes.NpgsqlDbType.Text, dto.Name);
        cmd.Parameters.AddWithValue("username", NpgsqlTypes.NpgsqlDbType.Text, dto.Username);
        cmd.Parameters.AddWithValue("password", NpgsqlTypes.NpgsqlDbType.Text, dto.Password);

        var insertedIdObj = await cmd.ExecuteScalarAsync();
        if (insertedIdObj == null || insertedIdObj == DBNull.Value)
        {
            return Results.Problem(detail: "Failed to create user.", statusCode: 500);
        }

        var insertedId = Convert.ToInt64(insertedIdObj);

        var result = new { id = insertedId, name = dto.Name, username = dto.Username };
        return Results.Created($"/users/{insertedId}", result);
    }
    catch (Exception ex)
    {
        app.Logger.LogError(ex, "Error while creating user with username {Username}", dto?.Username);
        return Results.Problem(detail: "An error occurred while creating the user.", statusCode: 500);
    }
}).WithName("createUser").WithOpenApi();

app.Run();

public record UserCreateDto(string Name, string Username, string Password);