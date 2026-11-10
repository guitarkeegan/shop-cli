import {
  onCleanup,
  createEffect,
  createMemo,
  on,
  onMount,
  createResource,
  createSignal,
  Accessor,
  Setter,
} from 'solid-js';
import { customElement, noShadowDOM } from 'solid-element';
import unique from 'lodash.uniq';
import './collection-product-grid.css';

type FilterValue = string | string[];
type Filters = Record<string, FilterValue>;

function formDataToObject(formData: FormData): Filters {
  const obj: Filters = {};
  for (const [key, value] of formData.entries()) {
    // Check if the key already exists
    if (obj.hasOwnProperty(key)) {
      // If it's not an array, convert it to an array
      if (!Array.isArray(obj[key])) {
        obj[key] = [obj[key] as string];
      }
      // Push the new value to the array
      (obj[key] as string[]).push(value as string);
    } else {
      // Assign the value directly
      obj[key] = value as string;
    }
  }
  return obj;
}

function outputToDebug(error: Error): void {
  const debugOutput = document.getElementById('debugOutput');
  if (!debugOutput) return;
  debugOutput.textContent = error.message;
}

function parseShopifyFiltersFromUrl(): Filters {
  const url = new URL(window.location.href);
  const filters: Filters = {};

  // Get all unique keys to handle repeated parameters
  const processedKeys = new Set<string>();

  url.searchParams.forEach((value, key) => {
    if (ignoredKeys.includes(key)) return;
    if (processedKeys.has(key)) return;
    processedKeys.add(key);

    // Handle Shopify filter params (filter.p.*, filter.v.*, etc.) directly
    if (key.startsWith('filter.')) {
      // Get all values for this key (handles repeated parameters for multi-select)
      const allValues = url.searchParams.getAll(key);
      filters[key] = allValues.length > 1 ? allValues : allValues[0];
    } else {
      // Regular params (sort_by, page, etc.)
      filters[key] = value;
    }
  });

  console.log('[parseShopifyFiltersFromUrl] Parsed filters from URL:', filters);
  return filters;
}

const createFilters = (element: HTMLElement, formSelector: Accessor<string>): [Accessor<Filters>, Setter<Filters>] => {
  const historyFilters = parseFiltersFromState(history.state?.filters);
  const urlFilters = parseShopifyFiltersFromUrl();

  console.log('[createFilters] historyFilters:', historyFilters);
  console.log('[createFilters] urlFilters:', urlFilters);

  const formFilters = createMemo(() => {
    const selector = formSelector();
    console.log('[createFilters] formFilters memo running, selector:', selector);
    const form = element.querySelector(selector) as HTMLFormElement | null;
    if (!form) {
      console.log('[createFilters] formFilters: form not found');
      return {};
    }
    const formData = new FormData(form);
    const filters = formDataToObject(formData);
    console.log('[createFilters] formFilters from FormData:', filters);
    return filters;
  });

  const mergedFilters = createMemo(() => {
    const merged = { ...urlFilters, ...historyFilters, ...formFilters() };
    console.log('[createFilters] mergedFilters:', merged);
    return merged;
  });

  const [filters, setFilters] = createSignal<Filters>(mergedFilters());

  onMount(() => {
    function handlePopState(event: PopStateEvent): void {
      if (event.state?.filters === JSON.stringify(filters())) return;
      setFilters(parseFiltersFromState(event.state?.filters));
    }
    addEventListener('popstate', handlePopState);

    return onCleanup(() => {
      removeEventListener('popstate', handlePopState);
    });
  });

  createEffect(
    on(mergedFilters, (mergedFilters) => {
      console.log('[createFilters effect] mergedFilters changed, updating filters signal:', mergedFilters);
      setFilters((filters) => ({ ...filters, ...mergedFilters }));
    }),
  );

  return [filters, setFilters];
};
function parseFiltersFromState(state: string | undefined): Filters {
  return state ? JSON.parse(decodeURIComponent(state)) : {};
}

type KeywordHandler = (url: URL, value?: FilterValue) => URL;

const keywordHandles: Record<string, KeywordHandler> = {
  tags: (url: URL, value?: FilterValue): URL => {
    const parts = url.pathname.split('/');
    let tagValue = value as string | string[] | undefined;
    if (Array.isArray(tagValue)) tagValue = tagValue.join('+');
    parts[3] = tagValue || '';
    if (!tagValue?.length) parts.pop();
    url.pathname = parts.join('/');
    return url;
  },
};

function getFormSelector(formId: Accessor<string | null> = () => null): string {
  return formId() ? `#${formId()}` : 'form';
}

const ignoredKeys = ['view', 'section_id'];

function getNewUrl(filters: Filters): URL {
  const url = new URL(window.location.href);
  const filterKeys = Object.keys(filters);

  // Clean up existing filter params that are no longer in filters
  Array.from(url.searchParams.keys()).forEach((key) => {
    if (ignoredKeys.includes(key)) return;

    // Check if this is a Shopify filter param (filter.p.*, filter.v.*, etc.)
    const isShopifyFilter = key.startsWith('filter.');

    if (isShopifyFilter) {
      // Remove if not in current filters
      if (!filterKeys.includes(key)) {
        url.searchParams.delete(key);
      }
    } else if (!filterKeys.includes(key)) {
      // Remove non-Shopify params that aren't in filters
      url.searchParams.delete(key);
    } else if (filterKeys.includes(key) && !(filters[key] as any).length) {
      // Remove if filter has no values
      url.searchParams.delete(key);
    }
  });

  filterKeys.forEach((key) => {
    let value = filters[key];

    // Handle special keyword handlers (e.g., tags for URL paths)
    if (keywordHandles[key]) {
      keywordHandles[key](url, value);
      return;
    }

    // Handle Shopify filter params directly (filter.p.*, filter.v.*, etc.)
    const isShopifyFilter = key.startsWith('filter.');
    if (isShopifyFilter) {
      const isArray = Array.isArray(value);

      if (!value || (isArray && !value.length)) {
        url.searchParams.delete(key);
        return;
      }

      // Shopify expects repeated parameters for multiple selections (OR logic)
      // e.g., filter.p.vendor=BioTE&filter.p.vendor=Cellcore
      url.searchParams.delete(key);
      if (isArray) {
        value.forEach((val) => {
          url.searchParams.append(key, val);
        });
      } else {
        url.searchParams.set(key, value as string);
      }
      return;
    }
    // Handle regular params (sort_by, page, etc.)
    else {
      const isArray = Array.isArray(value);
      if (!value || (isArray && !value.length)) return;

      if (isArray) {
        url.searchParams.delete(key);
        value.forEach((val) => {
          url.searchParams.append(key, val);
        });
      } else {
        url.searchParams.set(key, value as string);
      }
    }
  });

  return url;
}

interface CollectionFiltersHandlerProps {
  productsSelector: string;
  paginationSelector: string;
  filtersSelector: string;
  productCountSelector: string;
  formId: string;
  sectionId: string;
  listItemKey: string;
}

const CollectionFiltersHandler = (props: CollectionFiltersHandlerProps, { element }: { element: HTMLElement }) => {
  noShadowDOM();

  const domParser = new DOMParser();

  // Create a stable form selector accessor
  const formSelector = createMemo(() => getFormSelector(() => props.formId));
  const [filters, setFilters] = createFilters(element, formSelector);

  let controller: AbortController | null = null;
  let isInitialLoad = true;

  // Initialize previousCollectionTags from filters on mount
  const initialFilters = filters();
  let previousCollectionTags: string[] = Array.isArray(initialFilters['filter.p.tag'])
    ? initialFilters['filter.p.tag']
    : initialFilters['filter.p.tag']
      ? [initialFilters['filter.p.tag']]
      : [];
  console.log('[CollectionFiltersHandler] Initial collection tags:', previousCollectionTags);

  // Helper function to sync form inputs with filter state
  function syncFormWithFilters(): void {
    const currentFilters = filters();
    const form = element.querySelector(formSelector()) as HTMLFormElement | null;

    if (!form) return;

    Object.entries(currentFilters).forEach(([key, value]) => {
      const inputs = form.querySelectorAll(`[name="${key}"]`) as NodeListOf<HTMLInputElement>;

      inputs.forEach((input) => {
        if (input.type === 'checkbox') {
          // For checkboxes, check if the input's value is in the filter array
          if (Array.isArray(value)) {
            input.checked = value.includes(input.value);
          } else {
            input.checked = input.value === value;
          }
        } else if (input.type === 'radio') {
          // For radio buttons, check the one that matches the filter value
          input.checked = input.value === value;
        } else {
          // For other inputs, set the value directly
          input.value = value as string;
        }
      });
    });
    console.log('[syncFormWithFilters] Synced form inputs with filters:', currentFilters);
  }

  const [htmlText] = createResource(filters, async (filters): Promise<string> => {
    const url = getNewUrl(filters);

    // Skip fetch on initial load if we're already on the correct URL
    // This prevents unnecessary refetch and preserves initial checkbox states
    if (isInitialLoad) {
      isInitialLoad = false;
      const currentUrl = new URL(window.location.href);
      const targetUrl = new URL(url);
      currentUrl.searchParams.delete('section_id');
      targetUrl.searchParams.delete('section_id');

      if (currentUrl.toString() === targetUrl.toString()) {
        console.log('[createResource] Skipping initial fetch - already on correct URL');
        return '';
      }
    }

    url.searchParams.set('section_id', props.sectionId);

    if (controller) controller.abort();

    controller = new AbortController();

    try {
      console.log('[createResource] Fetching with filters:', filters, 'URL:', url.toString());
      const htmlText = await fetch(url, { signal: controller.signal }).then((res) => res.text());
      const jsonFilters = JSON.stringify(filters);
      url.searchParams.delete('section_id');
      window.history.replaceState({ filters: jsonFilters }, '', url);
      return htmlText;
    } catch (err) {
      console.warn(err);
      return '';
    }
  });

  createEffect(
    on(
      () => ({
        paginationSelector: props.paginationSelector,
        htmlText: htmlText(),
      }),
      ({ paginationSelector, htmlText }) => {
        if (!paginationSelector || !htmlText) return;

        const paginationTarget = element.querySelector(paginationSelector);
        if (!paginationTarget) return;

        const html = domParser.parseFromString(htmlText, 'text/html');
        const newListPagination = html.querySelector(paginationSelector);
        if (!newListPagination) return;

        paginationTarget.replaceChildren(...Array.from(newListPagination.children));
      },
    ),
  );

  createEffect(
    on(
      () => ({
        productsSelector: props.productsSelector,
        listItemKey: props.listItemKey,
        htmlText: htmlText(),
      }),
      ({ productsSelector, listItemKey, htmlText }) => {
        if (!productsSelector || !htmlText) return;

        const productsTarget = element.querySelector(productsSelector);
        if (!productsTarget) return;

        const html = domParser.parseFromString(htmlText, 'text/html');
        const newProducts = html.querySelector(productsSelector);
        if (!newProducts) return;

        const oldChildren = Array.from(productsTarget.children);
        const newChildren = Array.from(newProducts.children);
        const mergedChildren = newChildren.map((child) => {
          const existing = oldChildren.find(
            (oldChild) =>
              child.hasAttribute(listItemKey) &&
              oldChild.hasAttribute(listItemKey) &&
              oldChild.getAttribute(listItemKey) === child.getAttribute(listItemKey),
          );
          if (existing) return existing;
          return child;
        });
        productsTarget.replaceChildren(...mergedChildren);
      },
    ),
  );

  createEffect(
    on(
      () => ({
        filtersSelector: props.filtersSelector,
        htmlText: htmlText(),
      }),
      ({ filtersSelector, htmlText }) => {
        if (!filtersSelector || !htmlText) return;

        // Get the collection tags from current filters at the time the HTML arrives
        // This represents what collection tags the newly fetched HTML is for
        const currentCollectionTags = filters()['filter.p.tag'];
        const currentTagsArray: string[] = Array.isArray(currentCollectionTags)
          ? currentCollectionTags
          : currentCollectionTags
            ? [currentCollectionTags]
            : [];

        // Compare arrays to detect changes
        const collectionTagsChanged =
          previousCollectionTags.length !== currentTagsArray.length ||
          !previousCollectionTags.every((tag, idx) => currentTagsArray[idx] === tag);

        console.log('[filters effect]', {
          previousCollectionTags,
          currentTagsArray,
          collectionTagsChanged,
        });

        if (!collectionTagsChanged) {
          console.log('[filters effect] Skipping filter update - collection tags unchanged');
          return;
        }

        console.log(
          '[filters effect] Updating filter options. Previous:',
          previousCollectionTags,
          'New:',
          currentTagsArray,
        );
        previousCollectionTags = currentTagsArray;

        const filtersTarget = element.querySelector(filtersSelector);
        if (!filtersTarget) {
          console.warn('[filters effect] Could not find filters target:', filtersSelector);
          return;
        }

        const html = domParser.parseFromString(htmlText, 'text/html');
        const newFilters = html.querySelector(filtersSelector);
        if (!newFilters) {
          console.warn('[filters effect] Could not find filters in fetched HTML:', filtersSelector);
          return;
        }

        console.log('[filters effect] Replacing filter options children');
        filtersTarget.replaceChildren(...Array.from(newFilters.children));

        // Re-sync form inputs after DOM update
        syncFormWithFilters();
      },
    ),
  );

  createEffect(
    on(
      () => ({
        productCountSelector: props.productCountSelector,
        htmlText: htmlText(),
      }),
      ({ productCountSelector, htmlText }) => {
        if (!productCountSelector || !htmlText) return;

        const productCountTarget = element.querySelector(productCountSelector);
        if (!productCountTarget) return;

        const html = domParser.parseFromString(htmlText, 'text/html');
        const newProductCount = html.querySelector(productCountSelector);
        if (!newProductCount) return;

        productCountTarget.textContent = newProductCount.textContent;
      },
    ),
  );

  // Update mobile filter count badge
  createEffect(
    on(filters, (currentFilters) => {
      const mobileFiltersCount = document.getElementById('mobileFiltersCount');
      if (!mobileFiltersCount) return;

      // Count active filters (exclude non-filter params like page, sort_by)
      let count = 0;
      Object.entries(currentFilters).forEach(([key, value]) => {
        // Skip non-filter parameters
        if (ignoredKeys.includes(key) || key === 'page' || key === 'sort_by' || key === 'view') return;

        // Only count Shopify filter params (filter.*)
        if (key.startsWith('filter.')) {
          if (Array.isArray(value)) {
            count += value.length;
          } else if (value) {
            count += 1;
          }
        }
      });

      // Update badge visibility and text
      mobileFiltersCount.textContent = count.toString();
      mobileFiltersCount.style.display = count > 0 ? 'inline-flex' : 'none';
    }),
  );

  createEffect(
    on(
      () => ({ paginationSelector: props.paginationSelector }),
      ({ paginationSelector }) => {
        if (!paginationSelector) return;
        const paginationTarget = element.querySelector(paginationSelector);
        if (!paginationTarget) return;

        function handleClick(e: Event): void {
          e.preventDefault();
          try {
            let target: HTMLAnchorElement | null = null;
            if ((e.target as HTMLElement).tagName === 'A') target = e.target as HTMLAnchorElement;
            else target = (e.target as HTMLElement).closest('a');
            if (!target) return;
            const url = new URL(target.href);
            let newPage = url.searchParams.get('page');
            setFilters((filters) => ({ ...filters, page: newPage || '1' }));
          } catch (error) {
            outputToDebug(error as Error);
          }
        }
        paginationTarget.addEventListener('click', handleClick);
        return onCleanup(() => {
          paginationTarget.removeEventListener('click', handleClick);
        });
      },
    ),
  );

  onMount(() => {
    // Sync form inputs with initial filter state from URL
    syncFormWithFilters();

    function handleChange(e: Event): void {
      const target = e.target as HTMLInputElement;
      const type = target.type;
      switch (type) {
        case 'checkbox':
          return updateCheckboxFilter(target);
        case 'radio':
          setFilters((filters) => {
            const updated: Filters = {
              ...filters,
              [target.name]: target.value,
            };

            if (filters.page) updated.page = '1';
            return updated;
          });
          return;
        default:
          setFilters((filters) => {
            const updated: Filters = {
              ...filters,
              [target.name]: target.value,
            };
            if (filters.page) updated.page = '1';
            return updated;
          });
          return;
      }
    }
    function handleReset(e: Event): void {
      e.preventDefault();
      setTimeout(() => {
        const form = e.target as HTMLFormElement;
        const formObject = formDataToObject(new FormData(form));
        const url = new URL(window.location.href);
        url.search = '';
        setFilters({ tags: [], ...formObject });
        window.location.href = url.toString();
      });
    }

    function handleSubmit(e: Event): void {
      e.preventDefault();
      const formObject = formDataToObject(new FormData(e.target as HTMLFormElement));
      setFilters((filters) => ({ ...filters, ...formObject }));
    }

    element.addEventListener('input', handleChange);
    element.addEventListener('reset', handleReset);
    element.addEventListener('submit', handleSubmit);

    return onCleanup(() => {
      element.removeEventListener('input', handleChange);
      element.removeEventListener('reset', handleReset);
      element.removeEventListener('submit', handleSubmit);
    });
  });

  function updateCheckboxFilter(checkbox: HTMLInputElement): void {
    if (checkbox.checked) {
      setFilters((filters) => {
        const currentValue = filters[checkbox.name];
        const updated: Filters = { ...filters };

        if (!currentValue) {
          updated[checkbox.name] = [checkbox.value];
        } else if (!Array.isArray(currentValue)) {
          updated[checkbox.name] = [currentValue as string, checkbox.value];
        } else {
          updated[checkbox.name] = unique([...(filters[checkbox.name] as string[]), checkbox.value]);
        }

        if (filters.page) updated.page = '1';
        return updated;
      });
      return;
    }

    setFilters((filters) => {
      const currentValue = filters[checkbox.name];
      const updated: Filters = { ...filters };

      if (!Array.isArray(currentValue)) {
        updated[checkbox.name] = [];
      } else {
        updated[checkbox.name] = (filters[checkbox.name] as string[])?.filter((v) => v !== checkbox.value) || [];
      }

      if (filters.page) updated.page = '1';
      return updated;
    });
  }
};

customElement<CollectionFiltersHandlerProps>(
  'collection-filters-handler',
  {
    productsSelector: '',
    paginationSelector: '',
    filtersSelector: '',
    productCountSelector: '',
    formId: '',
    sectionId: '',
    listItemKey: 'id',
  },
  CollectionFiltersHandler,
);
