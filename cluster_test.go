package nutcracker

// Nutcracker
// Problem-based approach
// Clusterise (tests)
// Copyright © 2022-2024 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	"testing"

	"github.com/mash/gokmeans"
)

func TestClusterEasy(t *testing.T) {
	var observations []gokmeans.Node = []gokmeans.Node{
		gokmeans.Node{20.0, 20.0, 20.0, 20.0},
		gokmeans.Node{21.0, 21.0, 21.0, 21.0},
		gokmeans.Node{100.5, 100.5, 100.5, 100.5},
		gokmeans.Node{50.1, 50.1, 50.1, 50.1},
		gokmeans.Node{64.2, 64.2, 64.2, 64.2},
	}

	// Get a list of centroids and output the values
	if success, centroids := gokmeans.Train(observations, 3, 50); success {
		// Show the centroids
		fmt.Println("The centroids are")
		for _, centroid := range centroids {
			fmt.Println(centroid)
		}

		// Output the clusters
		fmt.Println("...")
		for _, observation := range observations {
			index := gokmeans.Nearest(observation, centroids)
			fmt.Println(observation, "belongs in cluster", index+1, ".")
		}
	} else {
		t.Errorf("Cluster not success")
	}
}

func TestClusterEasy2(t *testing.T) {
	var observations []gokmeans.Node = []gokmeans.Node{
		// left wall
		gokmeans.Node{0.0, 0.0},
		gokmeans.Node{0.0, 1.0},
		gokmeans.Node{0.0, 2.0},
		gokmeans.Node{0.0, 3.0},
		gokmeans.Node{0.0, 4.0},
		// right wall
		gokmeans.Node{4.0, 0.0},
		gokmeans.Node{4.0, 1.0},
		gokmeans.Node{4.0, 2.0},
		gokmeans.Node{4.0, 3.0},
		gokmeans.Node{4.0, 4.0},
		// ball
		gokmeans.Node{2.0, 4.0},
		// stuff
		gokmeans.Node{2.0, 0.0},
	}

	var centers []gokmeans.Node = []gokmeans.Node{
		// left wall
		gokmeans.Node{0.0, 2.0},
		// right wall
		gokmeans.Node{4.0, 2.0},
		// ball
		gokmeans.Node{2.0, 4.0},
		// stuff
		gokmeans.Node{2.0, 0.0},
	}

	// Get a list of centroids and output the values

	if success, centroids := gokmeans.Train(observations, 4, 50); success {
		// Show the centroids
		fmt.Println("The centroids are")
		for _, centroid := range centroids {
			fmt.Println(centroid)
		}

		// Output the clusters
		fmt.Println("...")
		for _, observation := range observations {
			index := gokmeans.Nearest(observation, centroids)
			fmt.Println(observation, "belongs in cluster", index+1, ".")
		}
	} else {
		t.Errorf("Cluster not success")
	}
	fmt.Println("------------------------------------------------")
	if success, centroids := gokmeans.Train2(observations, 4, 50, centers); success {
		// Show the centroids
		fmt.Println("The centroids are")
		for _, centroid := range centroids {
			fmt.Println(centroid)
		}

		// Output the clusters
		fmt.Println("...")
		for _, observation := range observations {
			index := gokmeans.Nearest(observation, centroids)
			fmt.Println(observation, "belongs in cluster", index+1, ".")
		}
	} else {
		t.Errorf("Cluster not success")
	}
}
