// SearchParams

export type SearchParams = {
  query?: string;
  pagination: {
    page: number;
    page_size: number;
  }
  order_key?: string;
  is_order_desc?: boolean;
  offer_type?: string;
  manufacturers?: string[];
  models?: string[];
  drives?: string[];
  colors?: string[];
  transmissions?: string[];
  fuel_types?: string[];
  price_range?: {
    min?: number;
    max?: number;
  };
  mileage_range?: {
    min?: number;
    max?: number;
  };
  year_range?: {
    min?: number;
    max?: number;
  };
  engine_power_range?: {
    min?: number;
    max?: number;
  };
  engine_capacity_range?: {
    min?: number;
    max?: number;
  };
  // offer_creation_date_range: {
  //   min: string | null;
  //   max: string | null;
  // };
}

export const parseIntOrUndefined = (
  value: string | undefined,
): number | undefined => {
  if (!value) return undefined;
  const parsed = parseInt(value, 10);
  return isNaN(parsed) ? undefined : parsed;
};

export const parseArrayOrUndefined = (
  value: string | string[] | undefined,
): string[] | undefined => {
  if (!value) return undefined;
  if (Array.isArray(value)) {
    return value.flatMap((v) => v.split(',').map(s => s.trim()).filter(Boolean));
  }
  return value.split(',').map(s => s.trim()).filter(Boolean);
};

export const trimAllAfterFirstSpace = (
  values: string[] | undefined,
): string[] | undefined => {
  if (!values) return undefined;
  return values.map(value => {
    const spaceIndex = value.indexOf(' ');
    if (spaceIndex === -1) return value;
    return value.substring(spaceIndex + 1);
  });
};