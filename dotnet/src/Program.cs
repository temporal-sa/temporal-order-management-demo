using DotNetOrderManagement;
using Microsoft.Extensions.Logging;
using Temporalio.Client;
using Temporalio.Common.EnvConfig;
using Temporalio.Worker;

var connectOptions = ClientEnvConfig.LoadClientConnectOptions();
connectOptions.LoggerFactory = LoggerFactory.Create(
    builder => builder.
        AddSimpleConsole(options => options.TimestampFormat = "[HH:mm:ss] ").
        SetMinimumLevel(LogLevel.Information));
var client = await TemporalClient.ConnectAsync(connectOptions);
Console.WriteLine("✅ Client connected to {0} in namespace '{1}'", connectOptions.TargetHost, connectOptions.Namespace);

using var tokenSource = new CancellationTokenSource();
Console.CancelKeyPress += (_, eventArgs) =>
{
    tokenSource.Cancel();
    eventArgs.Cancel = true;
};

var taskQueue = GetEnvVarWithDefault("TEMPORAL_TASK_QUEUE", "orders");

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
