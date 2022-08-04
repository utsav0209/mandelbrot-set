package main

import "math"

func abs(c complex128) float64 {
	r, i := real(c), imag(c)
	return math.Sqrt(r*r + i*i)
}

func mandelBrot(p complex128, maxIterations int) (iterations int) {
	z := 0 + 0i

	for ; iterations < maxIterations && abs(z) < 2.0; iterations++ {
		z = z*z + p
	}

	return
}
