/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React from 'react';
import SelectableButtonGroup from '../../../common/ui/SelectableButtonGroup.jsx';

const PricingCategories = ({ activeKey, setActiveKey, modelCategories, categoryCounts, availableCategories, t }) => {
  const items = Object.entries(modelCategories)
    .filter(([key]) => availableCategories.includes(key))
    .map(([key, category]) => ({
      value: key,
      label: category.label,
      icon: category.icon,
      tagCount: categoryCounts[key] || 0,
    }));

  return (
    <SelectableButtonGroup
      title={t('模型分类')}
      items={items}
      activeValue={activeKey}
      onChange={setActiveKey}
      t={t}
    />
  );
};

export default PricingCategories; 