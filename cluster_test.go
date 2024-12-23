package nutcracker

// Nutcracker
// Problem-based approach
// Clusterise (tests)
// Copyright © 2022-2024 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"reflect"
	"testing"

	"github.com/kelindar/dbscan"
	"github.com/lfritz/clustering"
	dbscan2 "github.com/lfritz/clustering/dbscan"
	"github.com/lfritz/clustering/index"
	"github.com/mash/gokmeans"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"github.com/stretchr/testify/assert"
)

/*
Может пригодиться библиотека https://github.com/lfritz/clustering
у неё есть ограничение - она работает с двумерными массивами
*/

func TestClusterEasy(t *testing.T) {
	var observations []gokmeans.Node = []gokmeans.Node{
		gokmeans.Node{20.0, 20.0, 20.0, 20.0},
		gokmeans.Node{21.0, 21.0, 21.0, 21.0},
		gokmeans.Node{100.5, 100.5, 100.5, 100.5},
		gokmeans.Node{50.1, 50.1, 50.1, 50.1},
		gokmeans.Node{64.2, 64.2, 64.2, 64.2},
	}

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

func TestClusterMeans(t *testing.T) {
	var d clusters.Observations
	for x := 0; x < 1024; x++ {
		d = append(d, clusters.Coordinates{
			rand.Float64(),
			rand.Float64(),
		})
	}

	// Partition the data points into 16 clusters
	km := kmeans.New()
	clusters, err := km.Partition(d, 16)
	if err != nil {
		t.Error(err)
	}
	for _, c := range clusters {
		fmt.Printf("Centered at x: %.2f y: %.2f\n", c.Center[0], c.Center[1])
		fmt.Printf("Matching data points: %+v\n\n", c.Observations)
	}
}

func TestClusterMeans2(t *testing.T) {
	var d clusters.Observations
	// left wall
	d = append(d, clusters.Coordinates{0.0, 0.0})
	d = append(d, clusters.Coordinates{0.0, 1.0})
	d = append(d, clusters.Coordinates{0.0, 2.0})
	d = append(d, clusters.Coordinates{0.0, 3.0})
	d = append(d, clusters.Coordinates{0.0, 4.0})
	// right wall
	d = append(d, clusters.Coordinates{4.0, 0.0})
	d = append(d, clusters.Coordinates{4.0, 1.0})
	d = append(d, clusters.Coordinates{4.0, 2.0})
	d = append(d, clusters.Coordinates{4.0, 3.0})
	d = append(d, clusters.Coordinates{4.0, 4.0})
	// ball
	d = append(d, clusters.Coordinates{2.0, 4.0})
	// stuff
	d = append(d, clusters.Coordinates{2.0, 0.0})

	// Partition the data points into 4 clusters
	km, _ := kmeans.NewWithOptions(0.5, nil) //New()
	clusters, err := km.Partition(d, 4)
	if err != nil {
		t.Error(err)
	}
	for _, c := range clusters {
		fmt.Printf("Centered at x: %.2f y: %.2f\n", c.Center[0], c.Center[1])
		fmt.Printf("Matching data points: %+v\n\n", c.Observations)
	}
}

func TestClusterDbscan(t *testing.T) {
	/*
		Кластер кластеризует точки с помощью метода DBSCAN.
		Для этого требуются два параметра: эпсилон и минимальное количество точек,
		необходимое для формирования плотной области (minDensity).
		Он начинается с произвольной отправной точки, которая еще не была посещена.
		Извлекается ε-окрестность этой точки, и если она содержит достаточно много точек, запускается кластер.
		В противном случае точка помечается как шум. Обратите внимание,
		что эта точка позже может быть найдена в ε-окружении достаточного размера другой точки и,
		следовательно, стать частью кластера.
	*/
	clusters := dbscan.Cluster(2, 1.0,
		SimplePoint{0.0, 1.0},
		SimplePoint{0.0, 0.5},
		SimplePoint{0.0, 0.0},
		SimplePoint{0.0, 5.0},
		SimplePoint{0.0, 4.5},
		SimplePoint{0.0, 4.0})

	assert.Equal(t, 2, len(clusters))
	if len(clusters) == 2 {
		assert.Equal(t, 3, len(clusters[0]))
		assert.Equal(t, 3, len(clusters[1]))
	}

	fmt.Println(clusters)
}

func TestClusterDbscanNoData(t *testing.T) {
	log.Println("Executing TestClusterNoData")

	clusters := dbscan.Cluster(3, 1.0)
	assert.Equal(t, 0, len(clusters))
}

func TestClusterDbscanPong(t *testing.T) {
	clusters := dbscan.Cluster(4, 8.0,
		// left wall
		SimplePoint{0.0, 0.0},
		SimplePoint{0.0, 1.0},
		SimplePoint{0.0, 2.0},
		SimplePoint{0.0, 3.0},
		SimplePoint{0.0, 4.0},
		SimplePoint{0.0, 5.0},
		SimplePoint{0.0, 6.0},
		// right wall
		SimplePoint{7.0, 0.0},
		SimplePoint{7.0, 1.0},
		SimplePoint{7.0, 2.0},
		SimplePoint{7.0, 3.0},
		SimplePoint{7.0, 4.0},
		SimplePoint{7.0, 5.0},
		SimplePoint{7.0, 6.0},
		// ball
		SimplePoint{3.0, 6.0},
		SimplePoint{4.0, 6.0},
		SimplePoint{3.0, 5.0},
		SimplePoint{4.0, 5.0},
		// stuff
		SimplePoint{3.0, 1.0},
		SimplePoint{4.0, 1.0},
		SimplePoint{3.0, 0.0},
		SimplePoint{4.0, 0.0})

	assert.Equal(t, 4, len(clusters))
	// if len(clusters) == 4 {
	// 	assert.Equal(t, 5, len(clusters[0]))
	// 	assert.Equal(t, 5, len(clusters[1]))
	// }

	// fmt.Println(clusters)
}

type SimplePoint struct {
	positionX float64
	positionY float64
}

func (s SimplePoint) DistanceTo(c dbscan.Point) float64 {
	in := c.(SimplePoint)

	dist := math.Hypot(s.positionX-in.positionX, s.positionY-in.positionY)
	// fmt.Printf("point1: %v point2: %v distance: %v\n", s, c, dist)
	return dist
}

func (s SimplePoint) Name() string {
	return fmt.Sprintf("%v:%v", s.positionX, s.positionY)
}

type Point interface {
	Name() string
	DistanceTo(Point) float64
	// GetX() float64
	// GetY() float64
}

// ============================== github.com/lfritz/clustering ============================

var testPoints = [][2]float64{
	// cluster a: 3 points
	{1, 8}, {1, 7}, {2, 7},
	// cluster b: 8 points
	{6, 8}, {7, 8},
	{5, 7}, {6, 7}, {7, 7}, {8, 7},
	{6, 6}, {7, 6},
	// cluster c: 8 points
	{2, 3}, {3, 3}, {1, 2}, {2, 2}, {3, 2}, {2, 1}, {3, 1},
	{4, 2}, // border point of both c and d
	// cluster d: 5 points
	{5, 3}, {5, 2}, {6, 2}, {5, 1},
}

var testIndex = index.NewTrivialIndex(testPoints)

func TestClusteringDbscan(t *testing.T) {
	expected := []int{
		clustering.Noise, clustering.Noise, clustering.Noise,
		0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1,
		2, 2, 2, 2, 2,
	}

	cl := dbscan2.Dbscan(testIndex, 1.1, 4)
	if !reflect.DeepEqual(cl, expected) {
		t.Errorf("Dbscan(testIndex, 1.1, 4)\nreturned: %v\nexpected: %v",
			cl, expected)
	}
}

var testPongPoints = [][2]float64{
	// left wall
	{0.0, 0.0},
	{0.0, 1.0},
	{0.0, 2.0},
	{0.0, 3.0},
	{0.0, 4.0},
	{0.0, 5.0},
	{0.0, 6.0},
	// right wall
	{7.0, 0.0},
	{7.0, 1.0},
	{7.0, 2.0},
	{7.0, 3.0},
	{7.0, 4.0},
	{7.0, 5.0},
	{7.0, 6.0},
	// ball
	{3.0, 6.0},
	{4.0, 6.0},
	{3.0, 5.0},
	{4.0, 5.0},
	// stuff
	{3.0, 1.0},
	{4.0, 1.0},
	{3.0, 0.0},
	{4.0, 0.0},
}

var indexPongPoints = index.NewTrivialIndex(testPongPoints)

func TestClusteringDbscanPong(t *testing.T) {
	expected := []int{
		0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3,
	}

	cl := dbscan2.Dbscan(indexPongPoints, 2.1, 4)
	fmt.Println(len(cl))
	fmt.Println(cl)
	if !reflect.DeepEqual(cl, expected) {
		t.Errorf("Dbscan(testIndex, 1.1, 4)\nreturned: %v\nexpected: %v",
			cl, expected)
	}
}

func TestPointsGroup(t *testing.T) {
	// TODO: implement me
}
