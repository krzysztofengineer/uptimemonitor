package test

import (
	"testing"
	"uptimemonitor"
)

func TestUptime(t *testing.T) {
	t.Run("uptime is empty when there are no checks", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		m, _ := tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "http://example.com",
		})

		tc.AssertEqual(m.Uptime, float32(0))
		tc.AssertEqual(m.AvgResponseTimeMs, int64(0))
	})

	t.Run("uptime is computed when check is created", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "http://example.com",
		})

		tc.Store.CreateCheck(t.Context(), uptimemonitor.Check{
			MonitorID:      1,
			StatusCode:     200,
			ResponseTimeMs: 100,
		})

		m, _ := tc.Store.GetMonitorByID(t.Context(), 1)

		tc.AssertEqual(m.N, int64(1))
		tc.AssertEqual(m.Uptime, float32(100))
		tc.AssertEqual(m.AvgResponseTimeMs, int64(100))
	})

	t.Run("it works with multiple checks", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "http://example.com",
		})

		tc.Store.CreateCheck(t.Context(), uptimemonitor.Check{
			MonitorID:      1,
			StatusCode:     200,
			ResponseTimeMs: 100,
		})

		tc.Store.CreateCheck(t.Context(), uptimemonitor.Check{
			MonitorID:      1,
			StatusCode:     200,
			ResponseTimeMs: 200,
		})

		m, _ := tc.Store.GetMonitorByID(t.Context(), 1)

		tc.AssertEqual(m.N, int64(2))
		tc.AssertEqual(m.Uptime, float32(100))
		tc.AssertEqual(m.AvgResponseTimeMs, int64(150))
	})

	t.Run("uptime is updated", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "http://example.com",
		})

		tc.Store.CreateCheck(t.Context(), uptimemonitor.Check{
			MonitorID:      1,
			StatusCode:     404,
			ResponseTimeMs: 100,
		})

		m, _ := tc.Store.GetMonitorByID(t.Context(), 1)

		tc.AssertEqual(m.N, int64(1))
		tc.AssertEqual(m.Uptime, float32(0))

		tc.Store.CreateCheck(t.Context(), uptimemonitor.Check{
			MonitorID:      1,
			StatusCode:     200,
			ResponseTimeMs: 100,
		})

		m, _ = tc.Store.GetMonitorByID(t.Context(), 1)

		tc.AssertEqual(m.N, int64(2))
		tc.AssertEqual(m.Uptime, float32(50))

		tc.Store.CreateCheck(t.Context(), uptimemonitor.Check{
			MonitorID:      1,
			StatusCode:     500,
			ResponseTimeMs: 100,
		})

		m, _ = tc.Store.GetMonitorByID(t.Context(), 1)

		tc.AssertEqual(m.N, int64(3))
		tc.AssertEqual(m.Uptime, float32(33.3))
	})
}
