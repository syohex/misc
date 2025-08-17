using System.Diagnostics;
using System.Runtime.InteropServices;
using System.Text;

namespace BlogGen;

public static class Clipboard
{
    public static async Task Copy(string text)
    {
        ProcessStartInfo? psi = null;
        if (OperatingSystem.IsWindows())
        {
            psi = GetProcessStartInfoForWindows();
        }
        else if (OperatingSystem.IsLinux())
        {
            if (IsWsl())
            {
                psi = GetProcessStartInfoForWsl();
            }
        }

        if (psi == null)
        {
            var os = RuntimeInformation.OSDescription;
            throw new Exception($"Unsupported Platform: {os}");
        }

        using var proc = new Process();
        proc.StartInfo = psi;
        proc.EnableRaisingEvents = false;
        proc.Start();

        await proc.StandardInput.WriteAsync(text);
        await proc.StandardInput.FlushAsync();
        proc.StandardInput.Close();

        var stdOutTask = proc.StandardOutput.ReadToEndAsync();
        await Task.WhenAll(stdOutTask, proc.WaitForExitAsync());
    }

    private static ProcessStartInfo GetProcessStartInfoForWindows()
    {
        var clipCommandPath = @"C:\Windows\System32\clip.exe";
        return new ProcessStartInfo
        {
            FileName = clipCommandPath,
            RedirectStandardInput = true,
            RedirectStandardOutput = true,
            UseShellExecute = false,
            CreateNoWindow = true,
            StandardOutputEncoding = Encoding.UTF8,
        };
    }

    private static ProcessStartInfo GetProcessStartInfoForWsl()
    {
        var clipCommandPath = "/mnt/c/Windows/System32/clip.exe";
        return new ProcessStartInfo
        {
            FileName = clipCommandPath,
            RedirectStandardInput = true,
            RedirectStandardOutput = true,
            UseShellExecute = false,
            CreateNoWindow = true,
            StandardOutputEncoding = Encoding.UTF8,
        };
    }

    private static bool IsWsl() => Environment.GetEnvironmentVariable("WSL_DISTRO_NAME") != null;
}