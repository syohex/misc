using BlogGen;

if (args.Length < 1)
{
    throw new Exception($"Usage: BlogGen <url> {args[0]}");
}

var homeDir = Environment.GetFolderPath(Environment.SpecialFolder.UserProfile);
if (string.IsNullOrEmpty(homeDir))
{
    throw new Exception("Cannot find home directory");
}

var configPath = Path.Join(homeDir, ".config", "blog", "config.yaml");
var config = await ConfigLoader.Load(configPath);
var url = args[0];

var parser = ParserFactory.Create(url);
var product = await parser.Parse(url, config);

Console.WriteLine($"Title = {product.Title}");
Console.WriteLine($"Image = {product.Image}");
Console.WriteLine($"Url = {product.Url}");

await Clipboard.Copy(product.Title);