var builder = WebApplication.CreateBuilder(args);

builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();

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