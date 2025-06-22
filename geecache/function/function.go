package function

func partition(s string) [][]string {
	n := len(s)
	var ans [][]string
	dp := make([][]bool, n)
	for i := range dp {
		dp[i] = make([]bool, n)
		for j := range dp[i] {
			dp[i][j] = true
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dp[i][j] = s[i] == s[j] && dp[i+1][j-1]
		}
	}
	var set []string
	var dfs func(idx int)
	dfs = func(idx int) {
		if idx == n {
			ans = append(ans, set)
			return
		}
		for j := idx; j < n; j++ {
			if dp[idx][j] {
				set = append(set, s[idx:j+1])
				dfs(j + 1)
				set = set[:len(set)-1]
			}
		}
	}
	dfs(0)
	return ans
}
