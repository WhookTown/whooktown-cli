import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { Box, Text } from 'ink';
export function TrafficPanel({ states }) {
    if (states.length === 0) {
        return (_jsx(Box, { children: _jsx(Text, { color: "gray", children: "No traffic states found" }) }));
    }
    return (_jsxs(Box, { flexDirection: "column", children: [_jsxs(Box, { children: [_jsx(Box, { width: 24, children: _jsx(Text, { bold: true, color: "white", children: "Layout ID" }) }), _jsx(Box, { width: 20, children: _jsx(Text, { bold: true, color: "white", children: "Density" }) }), _jsx(Box, { width: 10, children: _jsx(Text, { bold: true, color: "white", children: "Speed" }) }), _jsx(Box, { width: 10, children: _jsx(Text, { bold: true, color: "white", children: "Enabled" }) })] }), _jsx(Text, { color: "gray", children: '─'.repeat(64) }), states.map((state, i) => (_jsxs(Box, { children: [_jsx(Box, { width: 24, children: _jsx(Text, { children: truncate(state.layout_id || '-', 22) }) }), _jsx(Box, { width: 20, children: _jsx(DensityBar, { density: state.density ?? 0 }) }), _jsx(Box, { width: 10, children: _jsx(SpeedBadge, { speed: state.speed }) }), _jsx(Box, { width: 10, children: _jsx(EnabledBadge, { enabled: state.enabled ?? false }) })] }, state.layout_id || i))), _jsx(Box, { marginTop: 1, children: _jsxs(Text, { color: "gray", children: [states.length, " layout(s)"] }) }), _jsx(Box, { marginTop: 1, children: _jsx(Text, { color: "gray", dimColor: true, children: "Use \"wt traffic set <layoutId> --density <n>\" to change traffic" }) })] }));
}
function DensityBar({ density }) {
    const percent = Math.round(density);
    const barLength = 10;
    const filled = Math.round((percent / 100) * barLength);
    const empty = barLength - filled;
    return (_jsxs(Text, { children: [_jsx(Text, { color: "green", children: '█'.repeat(filled) }), _jsx(Text, { color: "gray", children: '░'.repeat(empty) }), _jsxs(Text, { color: "gray", children: [" ", percent.toString().padStart(3), "%"] })] }));
}
function SpeedBadge({ speed }) {
    const s = speed?.toLowerCase();
    switch (s) {
        case 'slow':
            return _jsx(Text, { color: "blue", children: "slow" });
        case 'normal':
            return _jsx(Text, { children: "normal" });
        case 'fast':
            return _jsx(Text, { color: "cyan", children: "fast" });
        default:
            return _jsx(Text, { color: "gray", children: "-" });
    }
}
function EnabledBadge({ enabled }) {
    return enabled ? (_jsx(Text, { color: "green", children: "on" })) : (_jsx(Text, { color: "gray", children: "off" }));
}
function truncate(str, maxLen) {
    if (str.length <= maxLen)
        return str;
    return str.slice(0, maxLen - 2) + '..';
}
//# sourceMappingURL=TrafficPanel.js.map