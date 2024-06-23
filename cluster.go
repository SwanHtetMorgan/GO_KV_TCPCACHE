package main

type ClusterMemoryDB struct {
	ClientID      int
	KeyValueCache *KeyValueCache
}

func NewCluster(connnection int, cache *KeyValueCache) *ClusterMemoryDB {
	return &ClusterMemoryDB{
		connnection,
		cache,
	}
}
func clusterSort(cluster []ClusterMemoryDB, low, high int) {
	if low < high {
		mid := (low + high) / 2
		clusterSort(cluster, low, mid)
		clusterSort(cluster, mid+1, high)
		merge(cluster, low, mid, high)
	}
}

// Merge function to merge two sorted subarrays
func merge(cluster []ClusterMemoryDB, low, mid, high int) {
	leftSize := mid - low + 1
	rightSize := high - mid

	leftCluster := make([]ClusterMemoryDB, leftSize)
	rightCluster := make([]ClusterMemoryDB, rightSize)

	for i := 0; i < leftSize; i++ {
		leftCluster[i] = cluster[low+i]
	}
	for j := 0; j < rightSize; j++ {
		rightCluster[j] = cluster[mid+1+j]
	}

	i, j := 0, 0
	k := low

	for i < leftSize && j < rightSize {
		if leftCluster[i].ClientID <= rightCluster[j].ClientID {
			cluster[k] = leftCluster[i]
			i++
		} else {
			cluster[k] = rightCluster[j]
			j++
		}
		k++
	}

	for i < leftSize {
		cluster[k] = leftCluster[i]
		i++
		k++
	}

	for j < rightSize {
		cluster[k] = rightCluster[j]
		j++
		k++
	}
}
