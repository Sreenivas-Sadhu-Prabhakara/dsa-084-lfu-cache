# LFU Cache

> Striver SDE Sheet 084 · Day 14 Stack & Queue II — a doubly linked list per frequency plus a running minimum frequency; eviction and promotion both O(1). Go, stdlib only, tested.

## Problem

A fixed-capacity cache evicting the **l**east **f**requently **u**sed key — and when frequencies tie, the least *recently* used among them. Get and Put in O(1). LRU's harder sibling.

## Approach

Group keys by exact use-count: one recency-ordered doubly linked list per frequency, found through a map. Each access "promotes" a node from its frequency list to the next one up (two O(1) splices). Eviction needs the lowest occupied frequency; a `minFreq` counter tracks it with two cheap facts — a fresh insert always resets it to 1, and promotion can only empty the current minimum's list, bumping it by one. The tie-break comes free: within a frequency list, the tail is the least recently touched. Three maps/lists cooperating, every move constant-time.

**Complexity:** Get/Put O(1); O(capacity) space.

## Run

    go test ./...   # verify
    go run .        # demo

---

One of 191 repos solving Striver's SDE Sheet in Go — stdlib only, one problem per repo, every solution tested. Browse the set: [`dsa-*` repos](https://github.com/Sreenivas-Sadhu-Prabhakara?tab=repositories&q=dsa-&type=source).
