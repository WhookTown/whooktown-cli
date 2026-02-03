import { Command } from 'commander';
import { createClient } from '../lib/client.js';
import { handleError } from '../lib/errors.js';
import { formatTable, formatJSON, truncate } from '../lib/output.js';

interface LayoutData {
  name?: string;
  buildings?: Array<{ id: string; name?: string; type?: string }>;
}

const listCommand = new Command('list')
  .description('List layouts')
  .option('-f, --format <format>', 'Output format: table, json', 'table')
  .option('-v, --verbose', 'Show building details')
  .action(async (options: { format: string; verbose?: boolean }) => {
    try {
      const client = createClient();
      const layouts = await client.ui.getLayouts();

      if (options.format === 'json') {
        console.log(formatJSON(layouts));
        return;
      }

      if (layouts.length === 0) {
        console.log('No layouts found');
        return;
      }

      // Table format
      const headers = ['ID', 'Name', 'Buildings', 'Received'];
      const rows = layouts.map((layout) => {
        const data = layout.data as LayoutData;
        const buildingCount = data?.buildings?.length || 0;
        return [
          truncate(layout.layout_id, 36),
          truncate(data?.name || '-', 20),
          String(buildingCount),
          layout.received_at ? new Date(layout.received_at).toLocaleString() : '-',
        ];
      });

      console.log(formatTable(headers, rows));
      console.log(`\n${layouts.length} layout(s)`);

      // Verbose mode: show buildings for each layout
      if (options.verbose) {
        console.log('\n--- Building Details ---\n');
        for (const layout of layouts) {
          const data = layout.data as LayoutData;
          console.log(`Layout: ${data?.name || layout.layout_id}`);

          if (data?.buildings && data.buildings.length > 0) {
            const buildingHeaders = ['ID', 'Name', 'Type'];
            const buildingRows = data.buildings.map((b) => [
              b.id,
              b.name || '-',
              b.type || '-',
            ]);
            console.log(formatTable(buildingHeaders, buildingRows));
          } else {
            console.log('  No buildings');
          }
          console.log('');
        }
      }
    } catch (err) {
      handleError(err);
    }
  });

export const layoutCommand = new Command('layout')
  .description('Manage layouts (read-only)')
  .addCommand(listCommand);
