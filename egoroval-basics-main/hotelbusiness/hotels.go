//go:build !solution

package hotelbusiness

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	maxD := 0
	for _, g := range guests {
		if g.CheckOutDate > maxD {
			maxD = g.CheckOutDate
		}
	}
	a := make([]int, maxD+1)
	for _, g := range guests {
		a[g.CheckInDate]++
		a[g.CheckOutDate]--
	}
	count := 0
	l := make([]Load, 0)
	for i, v := range a {
		if v != 0 {
			count += v
			l = append(l, Load{i, count})
		}
	}
	return l
}
