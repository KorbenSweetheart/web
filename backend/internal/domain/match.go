package domain

import (
	"math"
)

const (
	weightDist   = 0.25 // 25% spatial proximity
	weightAct    = 0.65 // 65% activity and skill alignment
	weightMode   = 0.10 // 10% interaction mode compatibility
	alphaAct     = 0.30 // preserves mutual high interest when skills diverge
	topActWeight = 0.60 // weight for primary anchor activity vs breadth

	minActivityScore = 0.15 // minimum activity similarity required to qualify
	minTotalScore    = 0.30 // minimum total match score required for recommendation
)

// CalculateMatchScore computes the compatibility score between a requesting profile and candidate profile.
// Returns the computed score and whether the candidate meets mutual radius, activity, and score thresholds.
func CalculateMatchScore(req, cand *Profile) (float64, bool) {
	distKm := haversineDistanceKm(req.Lat, req.Lon, cand.Lat, cand.Lon)
	if distKm > req.MaxRadius || distKm > cand.MaxRadius {
		return 0.0, false
	}

	reqMap := make(map[int64]ProfileActivity, len(req.Activities))
	for i := range req.Activities {
		reqMap[req.Activities[i].ActivityID] = req.Activities[i]
	}

	actScore := calculateMultiActivityScore(reqMap, cand.Activities)
	if actScore < minActivityScore {
		return 0.0, false
	}

	distScore := calculateDistanceScore(distKm, req.MaxRadius)
	modeScore := calculateInteractionModeScore(req.InteractionMode, cand.InteractionMode)

	totalScore := (weightDist * distScore) + (weightAct * actScore) + (weightMode * modeScore)
	if totalScore < minTotalScore {
		return totalScore, false
	}

	return totalScore, true
}

// IsCandidate evaluates whether two profiles are mutual match candidates.
// Validates profile completeness before evaluating match score.
func IsCandidate(p1, p2 *Profile) bool {
	if !IsProfileComplete(p1) || !IsProfileComplete(p2) {
		return false
	}

	_, isMatch := CalculateMatchScore(p1, p2)
	return isMatch
}

// IsProfileComplete checks that a profile has all required touchpoints filled in.
func IsProfileComplete(p *Profile) bool {
	if len(p.Activities) == 0 {
		return false
	}
	if p.Lat == 0 && p.Lon == 0 {
		return false
	}
	if p.InteractionMode == 0 {
		return false
	}

	return true
}

// calculateDistanceScore computes spatial proximity score using Gaussian decay.
func calculateDistanceScore(distKm, maxRadiusKm float64) float64 {
	if distKm <= 0 {
		return 1.0
	}
	if maxRadiusKm <= 0 || distKm >= maxRadiusKm {
		return 0.0
	}

	sigma := maxRadiusKm / 2.0
	return math.Exp(-0.5 * math.Pow(distKm/sigma, 2))
}

// haversineDistanceKm computes great-circle distance between two GPS coordinates in kilometers.
func haversineDistanceKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	lat1Rad := lat1 * (math.Pi / 180.0)
	lat2Rad := lat2 * (math.Pi / 180.0)

	sinDLat := math.Sin(dLat / 2.0)
	sinDLon := math.Sin(dLon / 2.0)

	a := sinDLat*sinDLat + math.Cos(lat1Rad)*math.Cos(lat2Rad)*sinDLon*sinDLon
	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1.0-a))

	return earthRadiusKm * c
}

// calculateMultiActivityScore computes multi-activity alignment using a two-tier weighted blend.
func calculateMultiActivityScore(reqMap map[int64]ProfileActivity, candActs []ProfileActivity) float64 {
	if len(reqMap) == 0 || len(candActs) == 0 {
		return 0.0
	}

	maxActScore := 0.0
	sumActScore := 0.0
	sharedCount := 0

	for i := range candActs {
		candAct := candActs[i]
		if reqAct, ok := reqMap[candAct.ActivityID]; ok {
			score := calculateActivityScore(reqAct, candAct)
			if score > maxActScore {
				maxActScore = score
			}
			sumActScore += score
			sharedCount++
		}
	}

	if sharedCount == 0 {
		return 0.0
	}

	unionCount := len(reqMap) + len(candActs) - sharedCount
	if unionCount <= 0 {
		unionCount = 1
	}

	actScore := topActWeight*maxActScore + (1.0-topActWeight)*(sumActScore/math.Sqrt(float64(unionCount)))
	if actScore > 1.0 {
		actScore = 1.0
	}

	return actScore
}

// calculateActivityScore computes normalized similarity for a single shared activity.
func calculateActivityScore(req, cand ProfileActivity) float64 {
	// Interest alignment: range [0.04, 1.0]
	interestScore := float64(int(req.InterestLevel)*int(cand.InterestLevel)) / 25.0

	// Experience penalty: range [0.0, 1.0]
	expDiff := math.Abs(float64(int(req.Experience) - int(cand.Experience)))
	expScore := 1.0 - (expDiff / 4.0)

	return interestScore * (alphaAct + (1.0-alphaAct)*expScore)
}

// calculateInteractionModeScore evaluates compatibility between interaction preferences.
func calculateInteractionModeScore(reqMode, candMode InteractionMode) float64 {
	if reqMode == OpenToAnything || candMode == OpenToAnything || reqMode == candMode {
		return 1.0
	}
	return 0.0
}
