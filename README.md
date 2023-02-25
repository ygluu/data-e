# DE创作背景(DE creative background)
Data Engine（数据引擎）用于独立模块间数据共享（Data Engine, For data sharing between independent modules）

在模块化开发当中要求模块逻辑代码具有高度的独立性，但模块数据又要共享给其它模块使用。如果开放模块接口给其它模块直接引用，就会因为互相引用依赖而降低模块的独立性。使用数据引擎（中间件）的公共接口来交换模块数据，就能避免引用依赖使得模块具有更高的独立性，提高模块的移植和移除效率，以及实现模块的热拔插等功能。本引擎在2018~2020间陆续创建定型，先后命名为：数据交换机、数据缓冲、数据中台等，最终定名：数据引擎。后续陆续调整目录结构和源码文件名及完善功能。
In modular development, module logic code is required to be highly independent, but module data must be shared with other modules. If the open module interface is directly referenced to other modules, the independence of the module will be reduced because of the cross-reference dependency. Using the common interface of the data engine (middleware) to exchange module data can avoid reference dependency and make the module more independent, improve the efficiency of module migration and removal, and realize the functions of module hot plugging. This engine was successively created and finalized between 2018 and 2020, and was successively named as data switch, data buffer, data center, etc., and finally named as data engine. The directory structure, source code file name and functions will be adjusted in succession.

此引擎集成：名字系统、数据仓库、条件系统、公式系统、运算集合、流程步骤控制系统、监听系统（数据变化监听和条件监听）、自定义函数系统、Go数学函数库、数据工作站。
This engine integrates: name system, data warehouse, condition system, formula system, calculation set, process step control system, monitoring system (data change monitoring and condition monitoring), user-defined function system, Go math function library, and data workstation.

# 用法(Usage)：
将源码放置于go项目的src/lib目录下，如：.../src/lib/data-e。
Place the source code in the src/lib directory of the go project, such .../src/lib/data-e。

# 软件著作开源协议及版权声明（Software copyright open source agreement and copyright statement）
1、使用权：依据原著的作者联系方式将使用者及使用版本信息告知原著作者后，任何组织或个人可免费永久获得使用权，此使用权的授权方式为最终使用者授权（即使用者无权再将使用权授权给其他组织或个人）。
Right of use: After informing the user and version information of the original author according to the contact information of the original author, any organization or individual can obtain the right of use for free and permanently. The way of authorization of the right of use is the authorization of the end user (that is, the user has no right to authorize the right of use to other organizations or individuals).

2、翻译权：任何组织或个人获得作者许可后方可翻译成其它程序语言的版本（下称：翻译版本）并附带原著的README.md文件，翻译版本适用使用权条款。
Translation right: Any organization or individual can translate into other program language versions (hereinafter referred to as the translated version) with the original README.md file after obtaining the permission of the author. The translated version is subject to the terms of the right of use.

3、改编权：任何组织或个人获得作者许可后方可改编成相同程序语言的版本（下称：改编版本）并附带原著的README.md文件，改编版本适用使用权条款。
Adaptation right: Any organization or individual can adapt to the version of the same program language (hereinafter referred to as the adaptation version) with the original README.md file after obtaining the permission of the author. The adaptation version is subject to the terms of the right of use.

4、署名权：翻译版本和改编版本的每个源码文件应该注明与原著一致的作者信息及联系方式及原著出处，可注明翻译者或改编者信息。
The right of authorship: each source code file of the translated version and the adapted version should indicate the author's information and contact information consistent with the original work, as well as the source of the original work. The information of the translator or the adapted person can be indicated.

5、版权等权益：原著、翻译版本、改编版本及使用前述版本编译成的二进制可执行版本的版权等法律许可的软件著作权益均归原著作者所有。
Copyright and other rights and interests: The copyright and other legally permitted software copyright rights and interests of the original work, translated version, adapted version and binary executable version compiled from the aforementioned version belong to the original author.

6、违约责任：任何违反本协议及声明的组织或个人和任何使用违反本协议及声明的软件著作的组织或个人，原著作者有权依法追究其法律责任、道德责任和经济责任（除追讨违法所得外，还包含但不仅限于版权费、经济损失费、赔偿金、利息（法律许可的上限）、名誉损失费、精神损失费、诉讼费、律师费、公证费、取证费、交通费、误工费、证人出具证明和出庭产生的所有费用）。
Liability for breach of contract: any organization or individual that violates this agreement and the statement and any organization or individual that uses the software works that violate this agreement and the statement, the original author has the right to pursue its legal, moral and economic responsibilities according to law (in addition to recovering the illegal income, it also includes but is not limited to the copyright fee, economic loss fee, compensation, interest (the upper limit permitted by law), reputation loss fee, moral loss fee, legal fee Attorney's fee, notarial fee, evidence collection fee, transportation fee, work delay fee, and all expenses incurred by the witness to issue the certificate and appear in court).

7、本协议及声明以保护作者权益为目的，未尽事宜在法律基础上仍然以保护作者权益为前提。
The purpose of this agreement and statement is to protect the rights and interests of the author. Matters not covered in this agreement and statement are still based on the premise of protecting the rights and interests of the author on a legal basis.

8、本协议及声明的条款由网络翻译工具翻译成中文之外的语系条款，如有语意歧义以中文条款语意为准。
The terms of this agreement and the statement are translated into language terms other than Chinese by online translation tools. In case of any ambiguity, the meaning of the Chinese terms shall prevail.

# 免责声明（Disclaimers）
1、作者不承担因使用此软件著作(含原著、翻译版本、改编版本及使用前述版本编译成的二进制可执行版本)而产生的任何法律责任、经济责任和道德责任。
he author does not assume any legal, economic and moral responsibilities arising from the use of this software work (including the original work, translated version, adapted version and the binary executable version compiled from the previous version).

2、作者没有任何责任和义务向任何组织或个人提供任何无偿的技术支持和服务。
The author has no responsibility or obligation to provide any free technical support and services to any organization or individual.

3、本软件著作并无主观故意抄袭等侵权行为，如有雷同纯属巧合，敬请及时告知以便第一时间删除侵权部分。
There is no infringement such as subjective intentional plagiarism in this software work. If there is any similarity, it is purely coincidental. Please inform us in time so as to delete the infringing part at the first time.

4、本声明的条款由网络翻译工具翻译成中文之外的语系条款，如有语意歧义以中文条款语意为准。
The terms of this statement are translated into language family terms other than Chinese by online translation tools. In case of semantic ambiguity, the meaning of the Chinese terms shall prevail.

# 作者信息(Author information)
Author: Yigui Lu (卢益贵/码客)
Contact WX/QQ: 48092788
Blog: https://blog.csdn.net/guestcode
