import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { Box, Text } from 'ink';
export function CameraPanel({ states }) {
    if (states.length === 0) {
        return (_jsx(Box, { children: _jsx(Text, { color: "gray", children: "No camera states found" }) }));
    }
    return (_jsxs(Box, { flexDirection: "column", children: [_jsxs(Box, { children: [_jsx(Box, { width: 24, children: _jsx(Text, { bold: true, color: "white", children: "Layout ID" }) }), _jsx(Box, { width: 12, children: _jsx(Text, { bold: true, color: "white", children: "Mode" }) }), _jsx(Box, { width: 15, children: _jsx(Text, { bold: true, color: "white", children: "Flyover Speed" }) })] }), _jsx(Text, { color: "gray", children: '─'.repeat(51) }), states.map((state, i) => (_jsxs(Box, { children: [_jsx(Box, { width: 24, children: _jsx(Text, { children: truncate(state.layout_id || '-', 22) }) }), _jsx(Box, { width: 12, children: _jsx(ModeBadge, { mode: state.mode }) }), _jsx(Box, { width: 15, children: _jsx(Text, { color: "gray", children: state.flyover_speed?.toFixed(1) || '-' }) })] }, state.layout_id || i))), _jsx(Box, { marginTop: 1, children: _jsxs(Text, { color: "gray", children: [states.length, " layout(s)"] }) }), _jsx(Box, { marginTop: 1, children: _jsx(Text, { color: "gray", dimColor: true, children: "Use \"wt camera set <layoutId> --mode <mode>\" to change camera mode" }) })] }));
}
function ModeBadge({ mode }) {
    const m = mode?.toLowerCase();
    switch (m) {
        case 'orbit':
            return _jsx(Text, { color: "cyan", children: "orbit" });
        case 'fps':
            return _jsx(Text, { color: "yellow", children: "fps" });
        case 'flyover':
            return _jsx(Text, { color: "magenta", children: "flyover" });
        default:
            return _jsx(Text, { color: "gray", children: "-" });
    }
}
function truncate(str, maxLen) {
    if (str.length <= maxLen)
        return str;
    return str.slice(0, maxLen - 2) + '..';
}
//# sourceMappingURL=CameraPanel.js.map