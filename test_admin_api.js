// 测试 admin API 的简单脚本
const axios = require('axios');

const BASE_URL = 'http://localhost:8080';

async function testAdminAPI() {
  try {
    console.log('测试 Admin API...');
    
    // 测试基础统计接口
    console.log('\n1. 测试基础统计接口:');
    const statsResponse = await axios.get(`${BASE_URL}/admin/stats`);
    console.log('状态码:', statsResponse.status);
    console.log('响应数据:', JSON.stringify(statsResponse.data, null, 2));
    
    // 测试详细统计接口
    console.log('\n2. 测试详细统计接口:');
    const detailedResponse = await axios.get(`${BASE_URL}/admin/stats/detailed`);
    console.log('状态码:', detailedResponse.status);
    console.log('响应数据:', JSON.stringify(detailedResponse.data, null, 2));
    
    // 测试增强版统计接口
    console.log('\n3. 测试增强版统计接口:');
    const enhancedResponse = await axios.get(`${BASE_URL}/admin/stats/enhanced`);
    console.log('状态码:', enhancedResponse.status);
    console.log('响应数据:', JSON.stringify(enhancedResponse.data, null, 2));
    
  } catch (error) {
    console.error('测试失败:', error.message);
    if (error.response) {
      console.error('响应状态:', error.response.status);
      console.error('响应数据:', error.response.data);
    }
  }
}

// 运行测试
testAdminAPI();
