int uniquePaths(int m, int n) {
  var dp = List.generate(m, (index) => List.filled(n, 0));
  dp[0][0] = 1;

  for (var i = 0; i < m; ++i) {
    for (var j = 0; j < n; ++j) {
      if (i != 0) {
        dp[i][j] += dp[i - 1][j];
      }
      if (j != 0) {
        dp[i][j] += dp[i][j - 1];
      }
    }
  }

  return dp[m - 1][n - 1];
}

void main() {
  assert(uniquePaths(3, 7) == 28);
  assert(uniquePaths(3, 2) == 3);
  assert(uniquePaths(3, 1) == 1);
  print("OK");
}
