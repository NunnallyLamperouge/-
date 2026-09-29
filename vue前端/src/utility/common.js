function parseDateTime(dateTimeString) {
    // 使用 JavaScript 的 Date 构造函数解析
    const date = new Date(dateTimeString);

    // 如果解析失败，返回一个错误消息
    if (isNaN(date.getTime())) {
        throw new Error("Invalid date string format");
    }

    // 提取年月日时分秒
    const year = date.getFullYear();
    const month = date.getMonth() + 1; // 月份从0开始，需要+1
    const day = date.getDate();
    const hours = date.getHours();
    const minutes = date.getMinutes();
    const seconds = date.getSeconds();

    // 返回结果对象
    return year + "年" +
    month + "月" +
    day + "日" +
    hours + "时" +
    minutes + "分" +
    seconds+ "秒";
}

export { parseDateTime };