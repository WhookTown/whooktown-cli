import React from 'react';
import { Box, Text } from 'ink';
import type { SensorState } from '@whooktown/sdk';
import type { SensorInfo } from '../App.js';

interface Props {
  sensors: SensorState[];
  sensorInfoMap: Map<string, SensorInfo>;
}

export function SensorPanel({ sensors, sensorInfoMap }: Props) {
  if (sensors.length === 0) {
    return (
      <Box>
        <Text color="gray">No sensors found</Text>
      </Box>
    );
  }

  return (
    <Box flexDirection="column">
      {/* Header */}
      <Box>
        <Box width={38}><Text bold color="white">ID</Text></Box>
        <Box width={16}><Text bold color="white">Name</Text></Box>
        <Box width={16}><Text bold color="white">Layout</Text></Box>
        <Box width={10}><Text bold color="white">Status</Text></Box>
        <Box width={10}><Text bold color="white">Activity</Text></Box>
        <Box width={12}><Text bold color="white">Updated</Text></Box>
      </Box>

      {/* Separator */}
      <Text color="gray">{'─'.repeat(102)}</Text>

      {/* Rows */}
      {sensors.map((sensor, i) => {
        const info = sensorInfoMap.get(sensor.id);
        return (
        <Box key={sensor.id || i}>
          <Box width={38}>
            <Text color="cyan">{sensor.id || '-'}</Text>
          </Box>
          <Box width={16}>
            <Text>{truncate(info?.buildingName || '-', 14)}</Text>
          </Box>
          <Box width={16}>
            <Text color="magenta">{truncate(info?.layoutName || '-', 14)}</Text>
          </Box>
          <Box width={10}>
            <StatusBadge status={sensor.data?.status as string} />
          </Box>
          <Box width={10}>
            <ActivityBadge activity={sensor.data?.activity as string} />
          </Box>
          <Box width={12}>
            <Text color="gray">
              {sensor.received_at
                ? new Date(sensor.received_at).toLocaleTimeString()
                : '-'}
            </Text>
          </Box>
        </Box>
        );
      })}

      {/* Footer */}
      <Box marginTop={1}>
        <Text color="gray">{sensors.length} sensor(s)</Text>
      </Box>
    </Box>
  );
}

function StatusBadge({ status }: { status?: string }) {
  const s = status?.toLowerCase();
  switch (s) {
    case 'online':
      return <Text color="green">online</Text>;
    case 'offline':
      return <Text color="gray">offline</Text>;
    case 'warning':
      return <Text color="yellow">warning</Text>;
    case 'critical':
      return <Text color="red">critical</Text>;
    default:
      return <Text color="gray">-</Text>;
  }
}

function ActivityBadge({ activity }: { activity?: string }) {
  const a = activity?.toLowerCase();
  switch (a) {
    case 'slow':
      return <Text color="blue">slow</Text>;
    case 'normal':
      return <Text>normal</Text>;
    case 'fast':
      return <Text color="cyan">fast</Text>;
    default:
      return <Text color="gray">-</Text>;
  }
}

function truncate(str: string, maxLen: number): string {
  if (str.length <= maxLen) return str;
  return str.slice(0, maxLen - 2) + '..';
}
