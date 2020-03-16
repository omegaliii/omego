package omego

import (
    "fmt"
    "strings"
)

type node struct {
    pattern  string // only left node would have pattern, others are ""
    part     string
    isWild   bool // exact match? true if part contains : or *
    children []*node
}

// Print the info of the node
func (n *node) String() string {
    return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// Recursively insert the node
func (n *node) insert(pattern string, parts []string, height int) {
    // Base Case: return when it is the last one of parts
    if len(parts) == height {
        n.pattern = pattern
        return
    }

    part := parts[height]
    child := n.matchChild(part) // find child which match the current part
    if child == nil {
        // If no such child found, create a new node and append to the children
        child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
        n.children = append(n.children, child)
    }
    // Recursive call
    child.insert(pattern, parts, height+1)
}

// Recursively search the node
func (n *node) search(parts []string, height int) *node {
    if len(parts) == height || strings.HasPrefix(n.part, "*") {
        if n.pattern == "" {
            // search failed
            return nil
        }
        return n
    }

    part := parts[height]
    children := n.matchChildren(part) // found out the child nodes which match the part

    // Iterate the match child nodes and recursively search 
    for _, child := range children {
        result := child.search(parts, height+1)
        if result != nil {
            return result
        }
    }

    return nil
}

//################################ 
//######## Helper Methods ######## 
//################################ 

// Find out the child node which match the part for the current node
// User for inserting
func (n *node) matchChild(part string) *node {
    for _, child := range n.children {
        if child.part == part || child.isWild {
            return child
        }
    }
    return nil
}

// Find out the children 
// Used for searching
func (n *node) matchChildren(part string) []*node {
    nodes := make([]*node, 0)
    for _, child := range n.children {
        if child.part == part || child.isWild {
            nodes = append(nodes, child)
        }
    }
    return nodes
}