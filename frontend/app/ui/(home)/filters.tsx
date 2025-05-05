"use client";

import { BaseFilterTemplate } from "@/app/ui/(home)/base-filter-template/base-filter-template";
import { fetchFilterFields } from "@/app/lib/data";
import { useEffect, useState } from "react";

export default function Filters() {
    const [filters, setFilters] = useState<{ name: string; options: string[] }[]>([]);

    useEffect(() => {
        async function fetchData() {
            const data = await fetchFilterFields();
            setFilters(data);
        }
        fetchData();
    }, []);

    function handleFilterChange(selected: string[]) {
        // TODO implement
    }

    return (
        <>
            <p className="px-2">Filters</p>
            {filters.map((filter, index) => (
                <BaseFilterTemplate
                    key={index}
                    name={filter.name}
                    options={filter.options}
                    onChange={handleFilterChange}
                />
            ))}
        </>
    );
}