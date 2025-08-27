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

        if (IsWsl())
        {
            await proc.StandardInput.BaseStream.WriteAsync(ConvertToUTF16(text));
        }
        else
        {
            await proc.StandardInput.WriteAsync(text);
        }

        await proc.StandardInput.FlushAsync();
        proc.StandardInput.Close();

        await proc.WaitForExitAsync();
    }

    private static byte[] ConvertToUTF16(string text) => Encoding.Unicode.GetBytes(text);

    private static ProcessStartInfo GetProcessStartInfoForWindows()
    {
        var clipCommandPath = @"C:\Windows\System32\clip.exe";
        return new ProcessStartInfo
        {
            FileName = clipCommandPath,
            RedirectStandardInput = true,
            UseShellExecute = false,
            CreateNoWindow = true,
        };
    }

    private static ProcessStartInfo GetProcessStartInfoForWsl()
    {
        var clipCommandPath = "/mnt/c/Windows/System32/clip.exe";
        return new ProcessStartInfo
        {
            FileName = clipCommandPath,
            RedirectStandardInput = true,
            UseShellExecute = false,
            CreateNoWindow = true,
        };
    }

    private static bool IsWsl() => Environment.GetEnvironmentVariable("WSL_DISTRO_NAME") != null;
}
