// Filters

export type FilterFieldData = {
  fieldName: string;
  options: string[];
  selected?: string[];
};

export type ModelFieldData = {
  producers: FilterFieldData;
  models: string[][];
}

// Ranges

export type RangeFieldData = {
  fieldName: string;
  range: {
    min: number | null;
    max: number | null;
  };
};
