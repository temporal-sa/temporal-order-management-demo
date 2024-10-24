using DotNetOrderManagement;
using Microsoft.Extensions.Logging;
using Temporalio.Client;
using Temporalio.Worker;

var address = GetEnvVarWithDefault("TEMPORAL_ADDRESS", "127.0.0.1:7233");
var temporalNamespace = GetEnvVarWithDefault("TEMPORAL_NAMESPACE", "default");
var tlsCertPath = GetEnvVarWithDefault("TEMPORAL_CERT_PATH", "");
var tlsKeyPath = GetEnvVarWithDefault("TEMPORAL_KEY_PATH", "");
var taskQueue = GetEnvVarWithDefault("TEMPORAL_TASK_QUEUE", "orders");
TlsOptions? tls = null;
if (!string.IsNullOrEmpty(tlsCertPath) && !string.IsNullOrEmpty(tlsKeyPath))
{
    Console.WriteLine("setting TLS");
    tls = new()
    {
        ClientCert = await File.ReadAllBytesAsync(tlsCertPath),
        ClientPrivateKey = await File.ReadAllBytesAsync(tlsKeyPath),
    };
}
Console.WriteLine($"Address is {address}");
var client = await TemporalClient.ConnectAsync(
    new(address)
    {
        Namespace = temporalNamespace,
        Tls = tls,
        LoggerFactory = LoggerFactory.Create(builder =>
            builder.
                AddSimpleConsole(options => options.TimestampFormat = "[HH:mm:ss] ").
                SetMinimumLevel(LogLevel.Information)),
    });

using var tokenSource = new CancellationTokenSource();
Console.CancelKeyPress += (_, eventArgs) =>
{
    tokenSource.Cancel();
    eventArgs.Cancel = true;
};

var activities = new OrderActivities();

using var worker = new TemporalWorker(
    client,
    new TemporalWorkerOptions(taskQueue).
        AddAllActivities(activities).
        AddWorkflow<OrderWorkflow>().
        AddWorkflow<OrderWorkflowScenarios>().
        AddWorkflow<ShippingChildWorkflow>());

// Run worker until cancelled
Console.WriteLine("Running worker...");
try
{
    await worker.ExecuteAsync(tokenSource.Token);
}
catch (OperationCanceledException)
{
    Console.WriteLine("Worker cancelled");
}

return 0;

string GetEnvVarWithDefault(string envName, string defaultValue)
{
    string? value = Environment.GetEnvironmentVariable(envName);
    if (string.IsNullOrEmpty(value))
    {
        return defaultValue;
    }
    return value;
}
