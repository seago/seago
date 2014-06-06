package utils

type Trie struct {
	Root *TrieNode
}

type TrieNode struct {
	Children map[rune]*TrieNode
	End      bool
}

func NewTrie() *Trie {
	t := &Trie{}
	t.Root = NewTrieNode()

	return t
}

func NewTrieNode() *TrieNode {
	n := &TrieNode{}
	n.Children = make(map[rune]*TrieNode)
	n.End = false

	return n
}

//新增要过滤的词
func (t *Trie) Add(txt string) {
	if len(txt) < 1 {
		return
	}
	chars := []rune(txt)
	slen := len(chars)
	node := t.Root
	for i := 0; i < slen; i++ {
		if _, exists := node.Children[chars[i]]; !exists {
			node.Children[chars[i]] = NewTrieNode()
		}
		node = node.Children[chars[i]]
	}
	node.End = true
}

//屏蔽字搜索查询
func (t *Trie) Find(txt string) bool {
	chars := []rune(txt)
	slen := len(chars)
	node := t.Root
	for i := 0; i < slen; i++ {
		if _, exists := node.Children[chars[i]]; exists {
			node = node.Children[chars[i]]
			for j := i + 1; j < slen; j++ {
				if _, exists := node.Children[chars[j]]; !exists {
					break
				}
				node = node.Children[chars[j]]
				if node.End == true {
					return true
				}
			}
			node = t.Root
		}
	}
	return false
}

//屏蔽字搜索替换
func (t *Trie) Replace(txt string) (string, []string) {
	chars := []rune(txt)
	result := []rune(txt)
	find := make([]string, 0, 10)
	slen := len(chars)
	node := t.Root
	for i := 0; i < slen; i++ {
		if _, exists := node.Children[chars[i]]; exists {
			node = node.Children[chars[i]]
			for j := i + 1; j < slen; j++ {
				if _, exists := node.Children[chars[j]]; !exists {
					break
				}
				node = node.Children[chars[j]]
				if node.End == true {
					for t := i; t <= j; t++ {
						result[t] = '*'
					}
					find = append(find, string(chars[i:j+1]))
					i = j
					node = t.Root
					break
				}
			}
			node = t.Root
		}
	}
	return string(result), find
}
