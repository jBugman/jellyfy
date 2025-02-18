package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

const version = "0.1.0"

type episode struct {
	fileName string

	number int
}

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	series := flag.String("series", "", "Series name")
	season := flag.Int("season", 0, "Season number")
	flag.Parse()

	if *showVersion {
		fmt.Printf("jellyfy version %s\n", version)
		return
	}

	folder := flag.Arg(0)

	if folder == "" || *series == "" || *season == 0 {
		fmt.Println("Usage: jellyfy -series=<series> -season=<season> <folder>")
		return
	}

	targetFolderName := fmt.Sprintf("Season %02d", *season)

	newFolderPath, err := renameFolder(folder, targetFolderName)
	if err != nil {
		fmt.Println("Error renaming folder:", err)
		return
	}

	episodes, err := listEpisodeFiles(newFolderPath)
	if err != nil {
		fmt.Println("Error listing episodes:", err)
		return
	}

	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].number < episodes[j].number
	})

	for _, episode := range episodes {
		ext := filepath.Ext(episode.fileName)
		targetName := fmt.Sprintf("%s S%02dE%02d%s", *series, *season, episode.number, ext)
		if targetName == episode.fileName {
			continue
		}

		srcPath := filepath.Join(newFolderPath, episode.fileName)
		newPath := filepath.Join(newFolderPath, targetName)
		if err := os.Rename(srcPath, newPath); err != nil {
			fmt.Println("Error renaming episode:", err)
		}
		fmt.Printf("%s -> %s\n", episode.fileName, targetName)
	}
}

func renameFolder(folderPath, newName string) (string, error) {
	parentDir := filepath.Dir(folderPath)
	newPath := filepath.Join(parentDir, newName)

	// Check if target already exists
	if _, err := os.Stat(newPath); err == nil {
		return "", fmt.Errorf("target folder already exists: '%s'", newName)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("error checking target folder: %w", err)
	}

	if err := os.Rename(folderPath, newPath); err != nil {
		return "", fmt.Errorf("error renaming folder: %w", err)
	}
	fmt.Println("Folder renamed successfully")

	return newPath, nil
}

func listEpisodeFiles(folderPath string) ([]episode, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	// Match patterns like S01E01, s1e1, S1E01, etc and capture episode number
	episodePattern := regexp.MustCompile(`(?i)s\d+e(\d+)`)
	var episodes []episode

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if matches := episodePattern.FindStringSubmatch(name); matches != nil {
			episodeNum, _ := strconv.Atoi(matches[1])
			episodes = append(episodes, episode{
				fileName: name,

				number: episodeNum,
			})
		}
	}

	return episodes, nil
}
