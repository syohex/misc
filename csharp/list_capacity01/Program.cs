var a = new List<int>(10);
Console.WriteLine($"length={a.Count} capacity={a.Capacity} a={string.Join(',', a)}");

for (var i = 0; i < 10; ++i) {
    a.Add(i);
}

Console.WriteLine($"length={a.Count} capacity={a.Capacity} a={string.Join(',', a)}");
