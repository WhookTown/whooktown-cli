import React from 'react';
import { Box, Text } from 'ink';
import type { TrafficState } from '@whooktown/sdk';

interface Props {
  states: TrafficState[];
}

export function TrafficPanel({ states }: Props) {
  if (states.length === 0) {
    return (
      <Box>
        <Text color="gray">No traffic states found</Text>
      </Box>
    );
  }

  return (
    <Box flexDirection="column">
      {/* Header */}
      <Box>
        <Box width={24}><Text bold color="white">Layout ID</Text></Box>
        <Box width={20}><Text bold color="white">Density</Text></Box>
        <Box width={10}><Text bold color="white">Speed</Text></Box>
        <Box width={10}><Text bold color="white">Enabled</Text></Box>
      </Box>

      {/* Separator */}
      <Text color="gray">{'─'.repeat(64)}</Text>

      {/* Rows */}
      {states.map((state, i) => (
        <Box key={state.layout_id || i}>
          <Box width={24}>
            <Text>{truncate(state.layout_id || '-', 22)}</Text>
          </Box>
          <Box width={20}>
            <DensityBar density={state.density ?? 0} />
          </Box>
          <Box width={10}>
            <SpeedBadge speed={state.speed} />
          </Box>
          <Box width={10}>
            <EnabledBadge enabled={state.enabled ?? false} />
          </Box>
        </Box>
      ))}

      {/* Footer */}
      <Box marginTop={1}>
        <Text color="gray">{states.length} layout(s)</Text>
      </Box>

      {/* Hint */}
      <Box marginTop={1}>
        <Text color="gray" dimColor>
          Use "wt traffic set &lt;layoutId&gt; --density &lt;n&gt;" to change traffic
        </Text>
      </Box>
    </Box>
  );
}

function DensityBar({ density }: { density: number }) {
  const percent = Math.round(density);
  const barLength = 10;
  const filled = Math.round((percent / 100) * barLength);
  const empty = barLength - filled;

  return (
    <Text>
      <Text color="green">{'█'.repeat(filled)}</Text>
      <Text color="gray">{'░'.repeat(empty)}</Text>
      <Text color="gray"> {percent.toString().padStart(3)}%</Text>
    </Text>
  );
}

function SpeedBadge({ speed }: { speed?: string }) {
  const s = speed?.toLowerCase();
  switch (s) {
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

function EnabledBadge({ enabled }: { enabled: boolean }) {
  return enabled ? (
    <Text color="green">on</Text>
  ) : (
    <Text color="gray">off</Text>
  );
}

function truncate(str: string, maxLen: number): string {
  if (str.length <= maxLen) return str;
  return str.slice(0, maxLen - 2) + '..';
}
