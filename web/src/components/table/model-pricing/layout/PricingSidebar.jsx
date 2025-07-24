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
import { Button } from '@douyinfe/semi-ui';
import PricingCategories from '../filter/PricingCategories';
import PricingGroups from '../filter/PricingGroups';
import PricingQuotaTypes from '../filter/PricingQuotaTypes';
import PricingEndpointTypes from '../filter/PricingEndpointTypes';
import PricingDisplaySettings from '../filter/PricingDisplaySettings';
import { resetPricingFilters } from '../../../../helpers/utils';

const PricingSidebar = ({
  showWithRecharge,
  setShowWithRecharge,
  currency,
  setCurrency,
  handleChange,
  setActiveKey,
  showRatio,
  setShowRatio,
  viewMode,
  setViewMode,
  filterGroup,
  setFilterGroup,
  filterQuotaType,
  setFilterQuotaType,
  filterEndpointType,
  setFilterEndpointType,
  currentPage,
  setCurrentPage,
  tokenUnit,
  setTokenUnit,
  loading,
  t,
  ...categoryProps
}) => {

  const handleResetFilters = () =>
    resetPricingFilters({
      handleChange,
      setActiveKey,
      availableCategories: categoryProps.availableCategories,
      setShowWithRecharge,
      setCurrency,
      setShowRatio,
      setViewMode,
      setFilterGroup,
      setFilterQuotaType,
      setFilterEndpointType,
      setCurrentPage,
      setTokenUnit,
    });

  return (
    <div className="p-4">
      <div className="flex items-center justify-between mb-6">
        <div className="text-lg font-semibold text-gray-800">
          {t('筛选')}
        </div>
        <Button
          theme="outline"
          type='tertiary'
          onClick={handleResetFilters}
          className="text-gray-500 hover:text-gray-700"
        >
          {t('重置')}
        </Button>
      </div>

      <PricingDisplaySettings
        showWithRecharge={showWithRecharge}
        setShowWithRecharge={setShowWithRecharge}
        currency={currency}
        setCurrency={setCurrency}
        showRatio={showRatio}
        setShowRatio={setShowRatio}
        viewMode={viewMode}
        setViewMode={setViewMode}
        tokenUnit={tokenUnit}
        setTokenUnit={setTokenUnit}
        loading={loading}
        t={t}
      />

      <PricingCategories
        {...categoryProps}
        setActiveKey={setActiveKey}
        loading={loading}
        t={t}
      />

      <PricingGroups
        filterGroup={filterGroup}
        setFilterGroup={setFilterGroup}
        usableGroup={categoryProps.usableGroup}
        groupRatio={categoryProps.groupRatio}
        models={categoryProps.models}
        loading={loading}
        t={t}
      />

      <PricingQuotaTypes
        filterQuotaType={filterQuotaType}
        setFilterQuotaType={setFilterQuotaType}
        models={categoryProps.models}
        loading={loading}
        t={t}
      />

      <PricingEndpointTypes
        filterEndpointType={filterEndpointType}
        setFilterEndpointType={setFilterEndpointType}
        models={categoryProps.models}
        loading={loading}
        t={t}
      />
    </div>
  );
};

export default PricingSidebar; 