var files = Directory.GetFiles("/usr/local/bin/");
Console.WriteLine($"files={string.Join("  ", files)}");

var dirs = Directory.GetDirectories("/usr/local");
Console.WriteLine($"files={string.Join("  ", dirs)}");

var f = Path.GetFileName(dirs[0]);
Console.WriteLine($"basename={f}");
