import { IUser } from "@/entities/user.entity";
import { ColumnDef, Row } from "@tanstack/vue-table";
import { h } from "vue";

export const columns: ColumnDef<IUser>[] = [
    {
        accessorKey: 'username',
        cell: ({ row }: {row: Row<IUser>}) => {
            const username: string = row.getValue('username')

            return h('div', {class: 'text-left font-medium'}, username)
        }
    },
]