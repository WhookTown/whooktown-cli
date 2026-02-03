import { Command } from 'commander';
import { createClient } from '../lib/client.js';
import { handleError } from '../lib/errors.js';
import { formatTable, formatJSON, truncate, success, error, } from '../lib/output.js';
import chalk from 'chalk';
/**
 * Parse comma-separated tags into array
 */
function parseTags(tagString) {
    return tagString
        .split(',')
        .map((t) => t.trim())
        .filter((t) => t.length > 0);
}
/**
 * Find a building by ID across all layouts
 */
async function findBuildingInLayouts(client, buildingId) {
    const layouts = await client.ui.getLayouts();
    for (const layoutDB of layouts) {
        const layoutData = layoutDB.data;
        if (layoutData?.buildings) {
            const building = layoutData.buildings.find((b) => b.id === buildingId);
            if (building) {
                return { layoutDB, layoutData, building };
            }
        }
    }
    error(`Building not found: ${buildingId}`);
    process.exit(1);
}
// Labels subcommand
const labelsCommand = new Command('labels')
    .description('Toggle building labels on/off')
    .argument('<layoutId>', 'Layout ID')
    .option('--on', 'Enable labels')
    .option('--off', 'Disable labels')
    .action(async (layoutId, options) => {
    try {
        if (options.on && options.off) {
            error('Cannot use both --on and --off');
            process.exit(1);
        }
        if (!options.on && !options.off) {
            error('Must specify --on or --off');
            process.exit(1);
        }
        const client = createClient();
        const enabled = options.on === true;
        await client.popup.toggleLabels(layoutId, enabled);
        success(`Labels ${enabled ? 'enabled' : 'disabled'} for layout`);
        console.log(chalk.gray(`  Layout: ${layoutId}`));
    }
    catch (err) {
        handleError(err);
    }
});
// Set subcommand
const setCommand = new Command('set')
    .description('Set building description/tags/notes')
    .argument('<buildingId>', 'Building ID (UUID)')
    .option('-d, --description <text>', 'Set description')
    .option('-t, --tags <tags>', 'Set comma-separated tags')
    .option('-n, --notes <text>', 'Set notes')
    .option('--clear-description', 'Clear description')
    .option('--clear-tags', 'Clear tags')
    .option('--clear-notes', 'Clear notes')
    .action(async (buildingId, options) => {
    try {
        const hasOption = options.description !== undefined ||
            options.tags !== undefined ||
            options.notes !== undefined ||
            options.clearDescription ||
            options.clearTags ||
            options.clearNotes;
        if (!hasOption) {
            error('At least one field option is required');
            console.error('Use: --description, --tags, --notes, --clear-description, --clear-tags, --clear-notes');
            process.exit(1);
        }
        const client = createClient();
        const { layoutDB, layoutData, building } = await findBuildingInLayouts(client, buildingId);
        // Apply modifications
        if (options.description !== undefined) {
            building.description = options.description;
        }
        if (options.clearDescription) {
            delete building.description;
        }
        if (options.tags !== undefined) {
            building.tags = parseTags(options.tags);
        }
        if (options.clearTags) {
            delete building.tags;
        }
        if (options.notes !== undefined) {
            building.notes = options.notes;
        }
        if (options.clearNotes) {
            delete building.notes;
        }
        // Reconstruct layout for update
        const updatePayload = {
            id: layoutDB.layout_id,
            name: layoutData.name || '',
            grid: layoutData.grid || { width: 10, height: 10 },
            buildings: layoutData.buildings || [],
            roads: layoutData.roads,
        };
        await client.ui.updateLayout(updatePayload);
        success(`Updated building: ${building.name || buildingId}`);
        if (options.description !== undefined) {
            console.log(chalk.gray(`  Description: ${options.description}`));
        }
        if (options.tags !== undefined) {
            console.log(chalk.gray(`  Tags: ${parseTags(options.tags).join(', ')}`));
        }
        if (options.notes !== undefined) {
            console.log(chalk.gray(`  Notes: ${options.notes}`));
        }
        if (options.clearDescription) {
            console.log(chalk.gray('  Description: cleared'));
        }
        if (options.clearTags) {
            console.log(chalk.gray('  Tags: cleared'));
        }
        if (options.clearNotes) {
            console.log(chalk.gray('  Notes: cleared'));
        }
    }
    catch (err) {
        handleError(err);
    }
});
// Get subcommand
const getCommand = new Command('get')
    .description('Get building metadata')
    .argument('<buildingId>', 'Building ID (UUID)')
    .option('-f, --format <format>', 'Output format: text, json', 'text')
    .action(async (buildingId, options) => {
    try {
        const client = createClient();
        const { layoutDB, layoutData, building } = await findBuildingInLayouts(client, buildingId);
        if (options.format === 'json') {
            console.log(formatJSON({
                id: building.id,
                name: building.name || null,
                type: building.type,
                layout_id: layoutDB.layout_id,
                layout_name: layoutData.name || null,
                description: building.description || null,
                tags: building.tags || [],
                notes: building.notes || null,
            }));
            return;
        }
        // Text format
        const label = (text) => chalk.bold.cyan(text);
        console.log(chalk.bold('Building'));
        console.log(`  ${label('ID:')}          ${building.id}`);
        console.log(`  ${label('Name:')}        ${building.name || chalk.gray('-')}`);
        console.log(`  ${label('Type:')}        ${building.type}`);
        console.log(`  ${label('Layout:')}      ${layoutData.name || layoutDB.layout_id}`);
        console.log('');
        console.log(chalk.bold('Metadata'));
        console.log(`  ${label('Description:')} ${building.description || chalk.gray('-')}`);
        console.log(`  ${label('Tags:')}        ${building.tags?.length ? building.tags.join(', ') : chalk.gray('-')}`);
        console.log(`  ${label('Notes:')}       ${building.notes || chalk.gray('-')}`);
    }
    catch (err) {
        handleError(err);
    }
});
// List subcommand
const listCommand = new Command('list')
    .description('List buildings with metadata')
    .argument('<layoutId>', 'Layout ID')
    .option('-f, --format <format>', 'Output format: table, json', 'table')
    .option('--tags <filter>', 'Filter by tags (comma-separated, any match)')
    .action(async (layoutId, options) => {
    try {
        const client = createClient();
        const layouts = await client.ui.getLayouts();
        const layoutDB = layouts.find((l) => l.layout_id === layoutId);
        if (!layoutDB) {
            error(`Layout not found: ${layoutId}`);
            process.exit(1);
        }
        const layoutData = layoutDB.data;
        let buildings = layoutData?.buildings || [];
        // Filter by tags if specified
        if (options.tags) {
            const filterTags = parseTags(options.tags).map((t) => t.toLowerCase());
            buildings = buildings.filter((b) => b.tags?.some((t) => filterTags.includes(t.toLowerCase())));
        }
        if (buildings.length === 0) {
            console.log(chalk.gray('No buildings found'));
            return;
        }
        if (options.format === 'json') {
            console.log(formatJSON(buildings.map((b) => ({
                id: b.id,
                name: b.name || null,
                type: b.type,
                description: b.description || null,
                tags: b.tags || [],
                notes: b.notes || null,
            }))));
            return;
        }
        // Table format
        const headers = ['ID', 'Name', 'Type', 'Tags', 'Description'];
        const rows = buildings.map((b) => [
            b.id,
            truncate(b.name || '-', 16),
            truncate(b.type, 14),
            truncate(b.tags?.join(', ') || '-', 16),
            truncate(b.description || '-', 16),
        ]);
        console.log(formatTable(headers, rows));
        console.log(`\n${buildings.length} building(s)`);
    }
    catch (err) {
        handleError(err);
    }
});
// Export main popup command
export const popupCommand = new Command('popup')
    .description('Manage popups and building metadata')
    .addCommand(labelsCommand)
    .addCommand(setCommand)
    .addCommand(getCommand)
    .addCommand(listCommand);
//# sourceMappingURL=popup.js.map