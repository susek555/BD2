import { FilterFieldData, RangeFieldData } from "@/app/lib/definitions/reviews";

export function syncFiltersWithParams(
  filters: FilterFieldData[],
  searchParams: URLSearchParams,
  ): FilterFieldData[] {
  return filters.map((filter) => {
    const fieldName = filter.fieldName.replace(/\s+/g, "");
    const selectedValues = searchParams.get(fieldName)?.split(",") ?? [];
    return { ...filter, selected: selectedValues };
  });
  }

  // Synchronizacja zakresów z URL
  export function syncRangesWithParams(
  ranges: RangeFieldData[],
  searchParams: URLSearchParams,
  ): RangeFieldData[] {
  return ranges.map((range) => {
    const fieldName = range.fieldName.replace(/\s+/g, "");
    const min = parseInt(searchParams.get(`${fieldName}_min`) ?? "0", 10);
    const max = parseInt(searchParams.get(`${fieldName}_max`) ?? "0", 10);
    return { ...range, range: { min, max } };
  });
  }

