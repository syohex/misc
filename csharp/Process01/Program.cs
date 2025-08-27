using System.Diagnostics;

const string UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36";

if (args.Length < 1)
{
    Console.Error.WriteLine("Usage: prog URL");
    return;
}

var processInfo = new ProcessStartInfo
{
    FileName = "curl",
    RedirectStandardOutput = true,
    UseShellExecute = false,
    CreateNoWindow = true,
};

var curlArgs = new[]
{
    "-A", UserAgent,
    "--tlsv1.2",
    "--tls-max", "1.2",
    "-b", "AGEAUTH=ok",
    "-L",
    "-s",
    args[0]
};

foreach (var arg in curlArgs)
{
    processInfo.ArgumentList.Add(arg);
}

var proc = new Process { StartInfo = processInfo };
proc.Start();

var output = await proc.StandardOutput.ReadToEndAsync();
await proc.WaitForExitAsync();

Console.WriteLine(output);
