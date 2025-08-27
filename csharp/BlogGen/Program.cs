using BlogGen;
using Scriban;

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

var template = """
<a href="{{ product.url }}" target="_blank">
<img src="{{ product.image }}" alt="{{ product.title }}" />
</a>

<p>
</p>
""";

var t = Template.Parse(template);
var result = t.Render(new { product });
Console.WriteLine(result);

await Clipboard.Copy(product.Title);
