import React, { useState, useEffect, useCallback } from 'react';
import { Box, Text, useInput, useApp } from 'ink';
import Spinner from 'ink-spinner';
import { createClient } from '../lib/client.js';
import type { SensorState, CameraState, TrafficState, LayoutDB } from '@whooktown/sdk';
import { SensorPanel } from './components/SensorPanel.js';
import { CameraPanel } from './components/CameraPanel.js';
import { TrafficPanel } from './components/TrafficPanel.js';

type Panel = 'sensors' | 'camera' | 'traffic';

interface LayoutData {
  name?: string;
  buildings?: Array<{ id: string; name?: string }>;
}

export interface SensorInfo {
  layoutName: string;
  buildingName: string;
}

interface Props {
  refreshInterval: number;
}

export function App({ refreshInterval }: Props) {
  const { exit } = useApp();
  const [activePanel, setActivePanel] = useState<Panel>('sensors');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastRefresh, setLastRefresh] = useState<Date | null>(null);

  // Data states
  const [sensors, setSensors] = useState<SensorState[]>([]);
  const [cameraStates, setCameraStates] = useState<CameraState[]>([]);
  const [trafficStates, setTrafficStates] = useState<TrafficState[]>([]);
  const [sensorInfoMap, setSensorInfoMap] = useState<Map<string, SensorInfo>>(new Map());

  const fetchData = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const client = createClient();

      const [sensorsData, cameraData, trafficData] = await Promise.all([
        client.sensors.getSensors(),
        client.sensors.getCameraStates(),
        client.sensors.getTrafficStates(),
      ]);

      // Try to fetch layouts
      let infoMap = new Map<string, SensorInfo>();
      try {
        const layouts = await client.ui.getLayouts();
        for (const layout of layouts) {
          const layoutData = layout.data as LayoutData;
          const layoutName = layoutData?.name || layout.layout_id;
          if (layoutData?.buildings) {
            for (const building of layoutData.buildings) {
              infoMap.set(building.id, {
                layoutName,
                buildingName: building.name || '-',
              });
            }
          }
        }
      } catch {
        // Silently ignore if we can't fetch layouts
      }

      setSensors(sensorsData);
      setCameraStates(cameraData);
      setTrafficStates(trafficData);
      setSensorInfoMap(infoMap);
      setLastRefresh(new Date());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch data');
    } finally {
      setLoading(false);
    }
  }, []);

  // Initial fetch
  useEffect(() => {
    fetchData();
  }, [fetchData]);

  // Auto-refresh
  useEffect(() => {
    const interval = setInterval(fetchData, refreshInterval);
    return () => clearInterval(interval);
  }, [fetchData, refreshInterval]);

  // Keyboard input
  useInput((input, key) => {
    if (input === 'q' || (key.ctrl && input === 'c')) {
      exit();
    }
    if (input === '1') setActivePanel('sensors');
    if (input === '2') setActivePanel('camera');
    if (input === '3') setActivePanel('traffic');
    if (input === 'r') {
      fetchData();
    }
  });

  return (
    <Box flexDirection="column" padding={1}>
      {/* Header */}
      <Box marginBottom={1}>
        <Text bold color="cyan">
          {' '}Whooktown CLI{' '}
        </Text>
        {loading && (
          <Text color="yellow">
            {' '}<Spinner type="dots" />{' '}
          </Text>
        )}
        {error && (
          <Text color="red"> {error}</Text>
        )}
      </Box>

      {/* Tab bar */}
      <Box marginBottom={1}>
        <Tab label="1 Sensors" active={activePanel === 'sensors'} />
        <Text> </Text>
        <Tab label="2 Camera" active={activePanel === 'camera'} />
        <Text> </Text>
        <Tab label="3 Traffic" active={activePanel === 'traffic'} />
      </Box>

      {/* Content */}
      <Box flexDirection="column" minHeight={10}>
        {activePanel === 'sensors' && <SensorPanel sensors={sensors} sensorInfoMap={sensorInfoMap} />}
        {activePanel === 'camera' && <CameraPanel states={cameraStates} />}
        {activePanel === 'traffic' && <TrafficPanel states={trafficStates} />}
      </Box>

      {/* Footer */}
      <Box marginTop={1} borderStyle="single" borderColor="gray" paddingX={1}>
        <Text color="gray">
          [1-3] Switch panels  [r] Refresh  [q] Quit
          {lastRefresh && (
            <Text>  |  Last: {lastRefresh.toLocaleTimeString()}</Text>
          )}
        </Text>
      </Box>
    </Box>
  );
}

function Tab({ label, active }: { label: string; active: boolean }) {
  return (
    <Text
      backgroundColor={active ? 'cyan' : undefined}
      color={active ? 'black' : 'gray'}
    >
      {' '}{label}{' '}
    </Text>
  );
}
