import React from 'react';
import { Box, Text } from 'ink';
import type { CameraState } from '@whooktown/sdk';

interface Props {
  states: CameraState[];
}

export function CameraPanel({ states }: Props) {
  if (states.length === 0) {
    return (
      <Box>
        <Text color="gray">No camera states found</Text>
      </Box>
    );
  }

  return (
    <Box flexDirection="column">
      {/* Header */}
      <Box>
        <Box width={24}><Text bold color="white">Layout ID</Text></Box>
        <Box width={12}><Text bold color="white">Mode</Text></Box>
        <Box width={15}><Text bold color="white">Flyover Speed</Text></Box>
      </Box>

      {/* Separator */}
      <Text color="gray">{'─'.repeat(51)}</Text>

      {/* Rows */}
      {states.map((state, i) => (
        <Box key={state.layout_id || i}>
          <Box width={24}>
            <Text>{truncate(state.layout_id || '-', 22)}</Text>
          </Box>
          <Box width={12}>
            <ModeBadge mode={state.mode} />
          </Box>
          <Box width={15}>
            <Text color="gray">
              {state.flyover_speed?.toFixed(1) || '-'}
            </Text>
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
          Use "wt camera set &lt;layoutId&gt; --mode &lt;mode&gt;" to change camera mode
        </Text>
      </Box>
    </Box>
  );
}

function ModeBadge({ mode }: { mode?: string }) {
  const m = mode?.toLowerCase();
  switch (m) {
    case 'orbit':
      return <Text color="cyan">orbit</Text>;
    case 'fps':
      return <Text color="yellow">fps</Text>;
    case 'flyover':
      return <Text color="magenta">flyover</Text>;
    default:
      return <Text color="gray">-</Text>;
  }
}

function truncate(str: string, maxLen: number): string {
  if (str.length <= maxLen) return str;
  return str.slice(0, maxLen - 2) + '..';
}
