'use client';

import { useState, useEffect } from 'react';
import { fetchHealth } from '@/lib/api';
import { SystemHealth } from '@/types';

export function useHealth() {
  const [health, setHealth] = useState<SystemHealth | null>(null);
  const [isOnline, setIsOnline] = useState<boolean>(false);
  const [latencyMs, setLatencyMs] = useState<number | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  const checkHealth = async () => {
    const t0 = performance.now();
    const res = await fetchHealth();
    const t1 = performance.now();

    if (res && res.status) {
      setHealth(res.status);
      setIsOnline(res.status.status === 'healthy');
      setLatencyMs(Math.round(t1 - t0));
    } else {
      setIsOnline(false);
      setLatencyMs(null);
    }
    setLoading(false);
  };

  useEffect(() => {
    checkHealth();
    const interval = setInterval(checkHealth, 15000);
    return () => clearInterval(interval);
  }, []);

  return { health, isOnline, latencyMs, loading, refetch: checkHealth };
}
