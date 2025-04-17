package pprofmcpagent

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/pprof/profile"
)

// ViewMode represents the type of view for profile data
type ViewMode string

const (
	ViewModeFlat  ViewMode = "flat"
	ViewModeCum   ViewMode = "cum"
	ViewModeGraph ViewMode = "graph"
)

func formatLocation(locs []*profile.Location) string {
	if len(locs) == 0 {
		return ""
	}

	loc := locs[0]
	if len(loc.Line) == 0 {
		return ""
	}

	line := loc.Line[0]
	if line.Function == nil {
		return ""
	}

	return fmt.Sprintf("%s:%d", line.Function.Name, line.Line)
}

// getTopSamples returns the profile data based on the specified view mode
func getTopSamples(p *profile.Profile, n int, viewMode ViewMode, profileType string) string {
	switch viewMode {
	case ViewModeCum:
		return getCumulativeView(p, n, profileType)
	case ViewModeGraph:
		return getGraphView(p, n, profileType)
	default: // ViewModeFlat
		return getFlatView(p, n, profileType)
	}
}

// aggregateSampleValues is a common function that aggregates profile sample values
func aggregateSampleValues(samples []*profile.Sample, getLocation func([]*profile.Location) string) map[string][]int64 {
	aggregated := make(map[string][]int64)

	for _, sample := range samples {
		loc := getLocation(sample.Location)
		if loc == "" {
			continue
		}

		if existing, ok := aggregated[loc]; ok {
			for i, v := range sample.Value {
				if i < len(existing) {
					existing[i] += v
				}
			}
		} else {
			valueCopy := make([]int64, len(sample.Value))
			copy(valueCopy, sample.Value)
			aggregated[loc] = valueCopy
		}
	}

	return aggregated
}

// getFlatView returns flat profile view (direct values for each location)
func getFlatView(p *profile.Profile, n int, profileType string) string {
	aggregatedSamples := aggregateSampleValues(p.Sample, formatLocation)
	return formatResults("Flat view (direct values)", aggregatedSamples, n, profileType)
}

// getCumulativeView returns cumulative profile view (including child functions)
func getCumulativeView(p *profile.Profile, n int, profileType string) string {
	getLocation := func(locs []*profile.Location) string {
		if len(locs) == 0 {
			return ""
		}
		return formatLocation([]*profile.Location{locs[0]})
	}

	aggregatedSamples := aggregateSampleValues(p.Sample, getLocation)
	return formatResults("Cumulative view (including children)", aggregatedSamples, n, profileType)
}

// getGraphView returns a call graph view of the profile
func getGraphView(p *profile.Profile, n int, profileType string) string {
	type nodeInfo struct {
		values   []int64
		children map[string][]int64
	}

	nodes := make(map[string]*nodeInfo)

	// Build the call graph
	for _, sample := range p.Sample {
		if len(sample.Location) == 0 {
			continue
		}

		// Process the call stack
		for i := 0; i < len(sample.Location); i++ {
			caller := formatLocation(sample.Location[i:])
			if caller == "" {
				continue
			}

			// Initialize or update node
			if _, exists := nodes[caller]; !exists {
				nodes[caller] = &nodeInfo{
					values:   make([]int64, len(sample.Value)),
					children: make(map[string][]int64),
				}
			}

			// Add values to the current node
			for j, v := range sample.Value {
				nodes[caller].values[j] += v
			}

			// Add values to child relationships
			if i+1 < len(sample.Location) {
				callee := formatLocation(sample.Location[i+1:])
				if callee == "" {
					continue
				}

				if _, exists := nodes[caller].children[callee]; !exists {
					nodes[caller].children[callee] = make([]int64, len(sample.Value))
				}
				for j, v := range sample.Value {
					nodes[caller].children[callee][j] += v
				}
			}
		}
	}

	// Sort nodes by total value
	type nodePair struct {
		name  string
		total int64
	}
	var sortedNodes []nodePair
	for name, info := range nodes {
		var total int64
		for _, v := range info.values {
			total += v
		}
		sortedNodes = append(sortedNodes, nodePair{name, total})
	}
	sort.Slice(sortedNodes, func(i, j int) bool {
		return sortedNodes[i].total > sortedNodes[j].total
	})

	// Build the output
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Call graph view (top %d nodes)\n", n))
	result.WriteString("Each node is followed by its children.\n\n")

	for i := 0; i < n && i < len(sortedNodes); i++ {
		nodeName := sortedNodes[i].name
		nodeInfo := nodes[nodeName]

		// Write node information
		result.WriteString(fmt.Sprintf("Node: %s\n", nodeName))
		result.WriteString(fmt.Sprintf("Values: %s\n", formatValues(nodeInfo.values, profileType)))

		// Sort and write children
		if len(nodeInfo.children) > 0 {
			result.WriteString("Children:\n")
			type childPair struct {
				name   string
				values []int64
			}
			var sortedChildren []childPair
			for childName, childValues := range nodeInfo.children {
				sortedChildren = append(sortedChildren, childPair{childName, childValues})
			}
			sort.Slice(sortedChildren, func(i, j int) bool {
				var totalI, totalJ int64
				for _, v := range sortedChildren[i].values {
					totalI += v
				}
				for _, v := range sortedChildren[j].values {
					totalJ += v
				}
				return totalI > totalJ
			})

			for _, child := range sortedChildren {
				result.WriteString(fmt.Sprintf("  %s: %s\n", child.name, formatValues(child.values, profileType)))
			}
		}
		result.WriteString("\n")
	}

	return result.String()
}

func formatResults(title string, samples map[string][]int64, n int, profileType string) string {
	type sampleInfo struct {
		location string
		value    []int64
	}

	// Convert map to slice for sorting
	var sampleSlice []sampleInfo
	for loc, val := range samples {
		sampleSlice = append(sampleSlice, sampleInfo{loc, val})
	}

	// Sort by the first value in descending order
	sort.Slice(sampleSlice, func(i, j int) bool {
		return sampleSlice[i].value[0] > sampleSlice[j].value[0]
	})

	// Build the output string
	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s (showing top %d locations)\n\n", title, n))

	// Show top N samples
	for i := 0; i < n && i < len(sampleSlice); i++ {
		sample := sampleSlice[i]
		result.WriteString(fmt.Sprintf("%s: %s\n", sample.location, formatValues(sample.value, profileType)))
	}

	return result.String()
}

func formatValues(values []int64, profileType string) string {
	if len(values) == 0 {
		return "no values"
	}

	var parts []string
	for i, v := range values {
		var formatted string
		switch profileType {
		case ProfileTypeHeap:
			switch i {
			case 0:
				formatted = fmt.Sprintf("%s in use", formatValue(v))
			case 1:
				formatted = fmt.Sprintf("%s total alloc", formatValue(v))
			default:
				formatted = formatValue(v)
			}
		case ProfileTypeBlock:
			switch i {
			case 0:
				formatted = fmt.Sprintf("%d contentions", v)
			case 1:
				formatted = fmt.Sprintf("%v delay", time.Duration(v))
			default:
				formatted = formatValue(v)
			}
		case ProfileTypeCPU:
			if i == 0 {
				formatted = fmt.Sprintf("%v CPU time", time.Duration(v))
			} else {
				formatted = formatValue(v)
			}
		default:
			formatted = formatValue(v)
		}
		parts = append(parts, formatted)
	}

	return strings.Join(parts, ", ")
}

func formatValue(v int64) string {
	const (
		_B  = int64(1)
		_KB = _B * 1024
		_MB = _KB * 1024
		_GB = _MB * 1024
		_TB = _GB * 1024
	)

	switch {
	case v > _TB:
		return fmt.Sprintf("%.2fTB", float64(v)/float64(_TB))
	case v > _GB:
		return fmt.Sprintf("%.2fGB", float64(v)/float64(_GB))
	case v > _MB:
		return fmt.Sprintf("%.2fMB", float64(v)/float64(_MB))
	case v > _KB:
		return fmt.Sprintf("%.2fKB", float64(v)/float64(_KB))
	default:
		return fmt.Sprintf("%dB", v)
	}
}
