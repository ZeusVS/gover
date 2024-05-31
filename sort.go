package main

import (
	"cmp"
	"os"
	"slices"
)

func sortDirsFirst(dirs []os.DirEntry) []os.DirEntry {
	// First we will sort the files alphabetically
	slices.SortFunc(dirs, func(dir1, dir2 os.DirEntry) int {
		return cmp.Compare(dir1.Name(), dir2.Name())
	})

	// Now we will sort the dirs first (stable to keep alphabetic sorting)
	slices.SortStableFunc(dirs, func(dir1, dir2 os.DirEntry) int {
		if dir1.IsDir() && !dir2.IsDir() {
			return -1
		}
		if !dir1.IsDir() && dir2.IsDir() {
			return 1
		}
		return 0
	})

	return dirs
}

func sortDirsLast(dirs []os.DirEntry) []os.DirEntry {
	// First we will sort the files alphabetically
	slices.SortFunc(dirs, func(dir1, dir2 os.DirEntry) int {
		return cmp.Compare(dir1.Name(), dir2.Name())
	})

	// Now we will sort the dirs last (stable to keep alphabetic sorting)
	slices.SortStableFunc(dirs, func(dir1, dir2 os.DirEntry) int {
		if dir1.IsDir() && !dir2.IsDir() {
			return 1
		}
		if !dir1.IsDir() && dir2.IsDir() {
			return -1
		}
		return 0
	})

	return dirs
}

func sortAlpha(dirs []os.DirEntry) []os.DirEntry {
	// Sort the files alphabetically
	slices.SortFunc(dirs, func(dir1, dir2 os.DirEntry) int {
		return cmp.Compare(dir1.Name(), dir2.Name())
	})

	return dirs
}

func sortAlphaReverse(dirs []os.DirEntry) []os.DirEntry {
	// Sort the files alphabetically reversed
	slices.SortFunc(dirs, func(dir1, dir2 os.DirEntry) int {
		return cmp.Compare(dir2.Name(), dir1.Name())
	})

	return dirs
}

func sortDateOldest(dirs []os.DirEntry) []os.DirEntry {
	// First we will sort the files alphabetically
	slices.SortFunc(dirs, func(dir1, dir2 os.DirEntry) int {
		return cmp.Compare(dir1.Name(), dir2.Name())
	})

	// Now we will sort by modification time, oldest first (stable to keep alphabetic sorting)
	slices.SortStableFunc(dirs, func(dir1, dir2 os.DirEntry) int {
		info1, err := dir1.Info()
		if err != nil {
			// If we can not get the info we can not sort so we return 0
			return 0
		}
		info2, err := dir2.Info()
		if err != nil {
			// If we can not get the info we can not sort so we return 0
			return 0
		}
		return cmp.Compare(info1.ModTime().Unix(), info2.ModTime().Unix())
	})

	return dirs
}

func sortDateNewest(dirs []os.DirEntry) []os.DirEntry {
	// First we will sort the files alphabetically
	slices.SortFunc(dirs, func(dir1, dir2 os.DirEntry) int {
		return cmp.Compare(dir1.Name(), dir2.Name())
	})

	// Now we will sort by modification time, newest first (stable to keep alphabetic sorting)
	slices.SortStableFunc(dirs, func(dir1, dir2 os.DirEntry) int {
		info1, err := dir1.Info()
		if err != nil {
			// If we can not get the info we can not sort so we return 0
			return 0
		}
		info2, err := dir2.Info()
		if err != nil {
			// If we can not get the info we can not sort so we return 0
			return 0
		}
		return cmp.Compare(info2.ModTime().Unix(), info1.ModTime().Unix())
	})

	return dirs
}

func sortSizeLargest(dirs []os.DirEntry) []os.DirEntry {
	// First we will sort the files alphabetically
	slices.SortFunc(dirs, func(dir1, dir2 os.DirEntry) int {
		return cmp.Compare(dir1.Name(), dir2.Name())
	})

	// Now we will sort by modification time, oldest first (stable to keep alphabetic sorting)
	slices.SortStableFunc(dirs, func(dir1, dir2 os.DirEntry) int {
		info1, err := dir1.Info()
		if err != nil {
			// If we can not get the info we can not sort so we return 0
			return 0
		}
		info2, err := dir2.Info()
		if err != nil {
			// If we can not get the info we can not sort so we return 0
			return 0
		}
		return cmp.Compare(info2.Size(), info1.Size())
	})

	return dirs
}

func sortSizeSmallest(dirs []os.DirEntry) []os.DirEntry {
	// First we will sort the files alphabetically
	slices.SortFunc(dirs, func(dir1, dir2 os.DirEntry) int {
		return cmp.Compare(dir1.Name(), dir2.Name())
	})

	// Now we will sort by modification time, oldest first (stable to keep alphabetic sorting)
	slices.SortStableFunc(dirs, func(dir1, dir2 os.DirEntry) int {
		info1, err := dir1.Info()
		if err != nil {
			// If we can not get the info we can not sort so we return 0
			return 0
		}
		info2, err := dir2.Info()
		if err != nil {
			// If we can not get the info we can not sort so we return 0
			return 0
		}
		return cmp.Compare(info1.Size(), info2.Size())
	})

	return dirs
}
