import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { Box, Text } from 'ink';
export function SensorPanel({ sensors, sensorInfoMap }) {
    if (sensors.length === 0) {
        return (_jsx(Box, { children: _jsx(Text, { color: "gray", children: "No sensors found" }) }));
    }
    return (_jsxs(Box, { flexDirection: "column", children: [_jsxs(Box, { children: [_jsx(Box, { width: 38, children: _jsx(Text, { bold: true, color: "white", children: "ID" }) }), _jsx(Box, { width: 16, children: _jsx(Text, { bold: true, color: "white", children: "Name" }) }), _jsx(Box, { width: 16, children: _jsx(Text, { bold: true, color: "white", children: "Layout" }) }), _jsx(Box, { width: 10, children: _jsx(Text, { bold: true, color: "white", children: "Status" }) }), _jsx(Box, { width: 10, children: _jsx(Text, { bold: true, color: "white", children: "Activity" }) }), _jsx(Box, { width: 12, children: _jsx(Text, { bold: true, color: "white", children: "Updated" }) })] }), _jsx(Text, { color: "gray", children: '─'.repeat(102) }), sensors.map((sensor, i) => {
                const info = sensorInfoMap.get(sensor.id);
                return (_jsxs(Box, { children: [_jsx(Box, { width: 38, children: _jsx(Text, { color: "cyan", children: sensor.id || '-' }) }), _jsx(Box, { width: 16, children: _jsx(Text, { children: truncate(info?.buildingName || '-', 14) }) }), _jsx(Box, { width: 16, children: _jsx(Text, { color: "magenta", children: truncate(info?.layoutName || '-', 14) }) }), _jsx(Box, { width: 10, children: _jsx(StatusBadge, { status: sensor.data?.status }) }), _jsx(Box, { width: 10, children: _jsx(ActivityBadge, { activity: sensor.data?.activity }) }), _jsx(Box, { width: 12, children: _jsx(Text, { color: "gray", children: sensor.received_at
                                    ? new Date(sensor.received_at).toLocaleTimeString()
                                    : '-' }) })] }, sensor.id || i));
            }), _jsx(Box, { marginTop: 1, children: _jsxs(Text, { color: "gray", children: [sensors.length, " sensor(s)"] }) })] }));
}
function StatusBadge({ status }) {
    const s = status?.toLowerCase();
    switch (s) {
        case 'online':
            return _jsx(Text, { color: "green", children: "online" });
        case 'offline':
            return _jsx(Text, { color: "gray", children: "offline" });
        case 'warning':
            return _jsx(Text, { color: "yellow", children: "warning" });
        case 'critical':
            return _jsx(Text, { color: "red", children: "critical" });
        default:
            return _jsx(Text, { color: "gray", children: "-" });
    }
}
function ActivityBadge({ activity }) {
    const a = activity?.toLowerCase();
    switch (a) {
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
function truncate(str, maxLen) {
    if (str.length <= maxLen)
        return str;
    return str.slice(0, maxLen - 2) + '..';
}
//# sourceMappingURL=SensorPanel.js.map