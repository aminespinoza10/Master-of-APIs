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

app.MapGet("/okCode", () =>
{
    return Results.Ok("Everything is awesome!");
}).WithName("getOkCode")
.WithOpenApi();

app.MapGet("/continueCode", () =>
{
    return Results.StatusCode(100);
}).WithName("getContinueCode")
.WithOpenApi();

app.MapGet("/movedPermanently", () =>
{
    return Results.StatusCode(301);
}).WithName("getMovedPermanently")
.WithOpenApi();

app.MapGet("/badRequest", () =>
{
    return Results.BadRequest("This was a bad request");
}).WithName("getBadRequest")
.WithOpenApi();

app.MapGet("/forbidden", () =>
{
    return Results.StatusCode(403);
}).WithName("getForbidden")
.WithOpenApi();

app.MapGet("/notFound", () =>
{
    return Results.NotFound("We couldn't find what you were looking for");
}).WithName("getNotFound")
.WithOpenApi();

app.MapGet("/proxyRequired", () =>
{
    return Results.StatusCode(407);
}).WithName("getProxyRequired")
.WithOpenApi();

app.Run();