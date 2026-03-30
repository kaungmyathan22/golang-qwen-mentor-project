// internal/algo/pathsum_test.go
package algo

import "testing"

func TestHasPathSum(t *testing.T) {
	tests := []struct {
		name      string
		root      *TreeNode
		targetSum int
		want      bool
	}{
		{
			name:      "valid path exists",
			root:      &TreeNode{Val: 5, Left: &TreeNode{Val: 4, Left: &TreeNode{Val: 11, Left: &TreeNode{Val: 7}, Right: &TreeNode{Val: 2}}}, Right: &TreeNode{Val: 8, Left: &TreeNode{Val: 13}, Right: &TreeNode{Val: 4, Right: &TreeNode{Val: 1}}}}, // build tree from Example 1
			targetSum: 22,
			want:      true,
		},
		{
			name:      "no valid path",
			root:      &TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 3}}, // build tree from Example 2
			targetSum: 5,
			want:      false,
		},
		{
			name:      "empty tree",
			root:      nil,
			targetSum: 0,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasPathSum(tt.root, tt.targetSum); got != tt.want {
				t.Errorf("HasPathSum() = %v, want %v", got, tt.want)
			}
		})
	}
}
