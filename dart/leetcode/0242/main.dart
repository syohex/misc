bool isAnagram(String s, String t) {
  if (s.length != t.length) {
    return false;
  }

  var sTable = List.filled(26, 0);
  var tTable = List.filled(26, 0);
  var baseIndex = "a".codeUnitAt(0);
  for (var i = 0; i < s.length; ++i) {
    sTable[s.codeUnitAt(i) - baseIndex]++;
    tTable[t.codeUnitAt(i) - baseIndex]++;
  }

  for (var i = 0; i < 26; ++i) {
    if (sTable[i] != tTable[i]) {
      return false;
    }
  }

  return true;
}

void main() {
  assert(isAnagram("margana", "anagram"));
  assert(!isAnagram("bomber", "mmomber"));
  print("OK");
}
