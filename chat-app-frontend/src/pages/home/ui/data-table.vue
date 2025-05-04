<script setup lang="ts" generic="TData, TValue">
    import { Button } from '@/components/ui/button';
    import { Table, TableBody, TableCell, TableRow } from '@/components/ui/table';
    import { ColumnDef, FlexRender, getCoreRowModel, useVueTable } from '@tanstack/vue-table';

    const props = defineProps<{
        columns: ColumnDef<TData, TValue>[],
        data: TData[]
    }>()

    const table = useVueTable({
        get data() { return props.data },
        get columns() { return props.columns },
        getCoreRowModel: getCoreRowModel(),
    })
</script>

<template>
    <div class="border rounded-md w-72">
        <Table>
            <TableBody>
                <template v-if="table.getRowModel().rows?.length">
                    <TableRow
                        v-for="row in table.getRowModel().rows" :key="row.id"
                        :data-state="row.getIsSelected() ? 'selected' : undefined"
                    >
                        <TableCell v-for="cell in row.getVisibleCells()" :key="cell.id">
                            <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
                        </TableCell>
                        <TableCell class="flex justify-center items-center">
                            <Button variant="outline" class="bg-blue-600 text-white">Написать</Button>
                        </TableCell>
                    </TableRow>
                </template>
                <template v-else>
                    <TableRow>
                        <TableCell :colspan="columns.length" class="h-24 text-center">
                            No Results
                        </TableCell>
                    </TableRow>
                </template>
            </TableBody>
        </Table>
    </div>
</template>