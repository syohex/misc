using YamlDotNet.Serialization;
using YamlDotNet.Serialization.NamingConventions;

namespace BlogGen;

public static class ConfigLoader
{
    public static async Task<Config> Load(string path)
    {
        var yaml = await File.ReadAllTextAsync(path);
        var deserializer = new DeserializerBuilder()
            .WithNamingConvention(CamelCaseNamingConvention.Instance)
            .IgnoreUnmatchedProperties()
            .Build();

        return deserializer.Deserialize<Config>(yaml);
    }
}