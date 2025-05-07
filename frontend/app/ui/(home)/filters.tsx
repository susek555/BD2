"use client";

import { BaseFilterTemplate } from "@/app/ui/(home)/base-filter-template/base-filter-template";
import { BaseRangeTemplate } from "./base-filter-template/base-range-template";
import { fetchFilterFields, prepareRangeFields } from "@/app/lib/data";
import { useEffect, useState } from "react";
import { FilterFieldData, RangeFieldData } from "@/app/lib/definitions";

export default function Filters() {
    const [filters, setFilters] = useState<FilterFieldData[]>([]);
    const [ranges, setRanges] = useState<RangeFieldData[]>([]);

    useEffect(() => {
        async function fetchData() {
            const data = await fetchFilterFields();
            setFilters(data);
            setRanges(prepareRangeFields());
        }
        fetchData();
    }, []);

    function handleFilterChange(name: string, selected: string[]) {
        // TODO implement
    }

    function handleRangeChange(name: string, range: { min: number; max: number }) {
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
            {ranges.map((range, index) => (
                <BaseRangeTemplate
                    key={index}
                    name={range.name}
                    onChange={handleRangeChange}
                />
            ))}
        </>
    );
}