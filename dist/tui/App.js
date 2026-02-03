import { jsxs as _jsxs, jsx as _jsx } from "react/jsx-runtime";
import { useState, useEffect, useCallback } from 'react';
import { Box, Text, useInput, useApp } from 'ink';
import Spinner from 'ink-spinner';
import { createClient } from '../lib/client.js';
import { SensorPanel } from './components/SensorPanel.js';
import { CameraPanel } from './components/CameraPanel.js';
import { TrafficPanel } from './components/TrafficPanel.js';
export function App({ refreshInterval }) {
    const { exit } = useApp();
    const [activePanel, setActivePanel] = useState('sensors');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [lastRefresh, setLastRefresh] = useState(null);
    // Data states
    const [sensors, setSensors] = useState([]);
    const [cameraStates, setCameraStates] = useState([]);
    const [trafficStates, setTrafficStates] = useState([]);
    const [sensorInfoMap, setSensorInfoMap] = useState(new Map());
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
            let infoMap = new Map();
            try {
                const layouts = await client.ui.getLayouts();
                for (const layout of layouts) {
                    const layoutData = layout.data;
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
            }
            catch {
                // Silently ignore if we can't fetch layouts
            }
            setSensors(sensorsData);
            setCameraStates(cameraData);
            setTrafficStates(trafficData);
            setSensorInfoMap(infoMap);
            setLastRefresh(new Date());
        }
        catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to fetch data');
        }
        finally {
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
        if (input === '1')
            setActivePanel('sensors');
        if (input === '2')
            setActivePanel('camera');
        if (input === '3')
            setActivePanel('traffic');
        if (input === 'r') {
            fetchData();
        }
    });
    return (_jsxs(Box, { flexDirection: "column", padding: 1, children: [_jsxs(Box, { marginBottom: 1, children: [_jsxs(Text, { bold: true, color: "cyan", children: [' ', "Whooktown CLI", ' '] }), loading && (_jsxs(Text, { color: "yellow", children: [' ', _jsx(Spinner, { type: "dots" }), ' '] })), error && (_jsxs(Text, { color: "red", children: [" ", error] }))] }), _jsxs(Box, { marginBottom: 1, children: [_jsx(Tab, { label: "1 Sensors", active: activePanel === 'sensors' }), _jsx(Text, { children: " " }), _jsx(Tab, { label: "2 Camera", active: activePanel === 'camera' }), _jsx(Text, { children: " " }), _jsx(Tab, { label: "3 Traffic", active: activePanel === 'traffic' })] }), _jsxs(Box, { flexDirection: "column", minHeight: 10, children: [activePanel === 'sensors' && _jsx(SensorPanel, { sensors: sensors, sensorInfoMap: sensorInfoMap }), activePanel === 'camera' && _jsx(CameraPanel, { states: cameraStates }), activePanel === 'traffic' && _jsx(TrafficPanel, { states: trafficStates })] }), _jsx(Box, { marginTop: 1, borderStyle: "single", borderColor: "gray", paddingX: 1, children: _jsxs(Text, { color: "gray", children: ["[1-3] Switch panels  [r] Refresh  [q] Quit", lastRefresh && (_jsxs(Text, { children: ["  |  Last: ", lastRefresh.toLocaleTimeString()] }))] }) })] }));
}
function Tab({ label, active }) {
    return (_jsxs(Text, { backgroundColor: active ? 'cyan' : undefined, color: active ? 'black' : 'gray', children: [' ', label, ' '] }));
}
//# sourceMappingURL=App.js.map