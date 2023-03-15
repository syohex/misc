using System.Text.RegularExpressions;

var timeRegex = new Regex(@"^(?:(?<hours>[0-9]+):)?(?<minutes>[0-9]{2}):(?<seconds>[0-9]{2}).(?<frac>[0-9]{3})$");
string[] input = {
    "00:09:22.000",
    "11:22.000",
};

foreach (var d in input)
{
    var m = timeRegex.Match(d);
    if (m.Success)
    {
        int hours = -1;
        int.TryParse(m.Groups["hours"].Value, out hours);
        var minutes = int.Parse(m.Groups["minutes"].Value);
        var seconds = int.Parse(m.Groups["seconds"].Value);
        var secondsFrac = int.Parse(m.Groups["frac"].Value);
        Console.WriteLine($"matched: hours='{hours}' minutes='{minutes}', seconds='{seconds}', frac='{secondsFrac}'");
    }
    else
    {
        Console.WriteLine($"{d} is not matched");
    }
}
