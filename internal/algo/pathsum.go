package algo

// internal/algo/pathsum.go

// TreeNode definition for binary tree
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// HasPathSum returns true if there exists a root-to-leaf path
// where the sum of values equals targetSum.
// Time: O(n), Space: O(h) where h = tree height
func HasPathSum(root *TreeNode, targetSum int) bool {
	if root == nil {
		return false
	}

	remaining := targetSum - root.Val

	if root.Left == nil && root.Right == nil {
		return remaining == 0
	}

	return HasPathSum(root.Left, remaining) || HasPathSum(root.Right, remaining)
}
