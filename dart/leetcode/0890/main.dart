List<int> toIndexPattern(String str) {
  List<int> ret = [];
  Map<int, int> m = new Map();

  for (var i = 0; i < str.length; ++i) {
    var c = str.codeUnitAt(i);
    if (m.containsKey(c)) {
      ret.add(m[c]!);
    } else {
      m[c] = i;
      ret.add(i);
    }
  }

  return ret;
}

bool equalLists<T>(List<T> a, List<T> b) {
  if (a.length != b.length) {
    return false;
  }

  for (var i = 0; i < a.length; ++i) {
    if (a[i] != b[i]) {
      return false;
    }
  }

  return true;
}

List<String> findAndReplacePattern(List<String> words, String pattern) {
  List<String> ret = [];

  final p = toIndexPattern(pattern);
  for (final word in words) {
    final w = toIndexPattern(word);
    if (equalLists(p, w)) {
      ret.add(word);
    }
  }

  return ret;
}

void main() {
  {
    final words = ["abc", "deq", "mee", "aqq", "dkd", "ccc"];
    final pattern = "abb";
    final expected = ["mee", "aqq"];

    final ret = findAndReplacePattern(words, pattern);
    assert(equalLists(ret, expected));
  }
  {
    final words = ["a", "b", "c"];
    final pattern = "a";
    final expected = ["a", "b", "c"];

    final ret = findAndReplacePattern(words, pattern);
    assert(equalLists(ret, expected));
  }

  print("OK");
}
